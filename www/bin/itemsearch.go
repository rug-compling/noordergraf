package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/http/cgi"
	"net/url"
	"runtime"
	"strings"
)

type Sparql struct {
	Results []ResultT `xml:"results>result"`
}

type ResultT struct {
	Bindings []BindingT `xml:"binding"`
}

type BindingT struct {
	Name    string `xml:"name,attr"`
	URI     string `xml:"uri"`
	Literal string `xml:"literal"`
}

func main() {
	req, err := cgi.Request()
	if x(err, http.StatusInternalServerError) {
		return
	}
	if x(req.ParseForm(), http.StatusInternalServerError) {
		return
	}
	q := req.FormValue("q")

	var query, title string

	switch q {
	default:
		fmt.Print(`Status: 400

Missing query
`)
		return
	case "famous":
		title = "beroemde mensen"
		query = `
SELECT ?plot ?name ?sameas {
  ?person :sameAs ?sameas .
  ?person :name / :fullName ?name .
  ?plot :subject ?person .
  FILTER fn:not(STRSTARTS(xsd:string(?sameas), "https://noordergraf"))
} ORDER BY ?name
`
	case "multi":
		title = "genoemd op meerdere graven"
		query = `
SELECT ?plot ?name ?sameas {
  ?person :sameAs ?sameas .
  ?person :name / :fullName ?name .
  ?plot :subject ?person .
  FILTER STRSTARTS(xsd:string(?sameas), "https://noordergraf")
} ORDER BY ?name
`
	}

	request := "http://localhost:10035/repositories/noordergraf?limit=1000&query=" + url.QueryEscape(query)
	fmt.Println(request)
	resp, err := http.Get(request)
	if x(err, http.StatusInternalServerError) {
		return
	}
	if resp.StatusCode > 299 {
		fmt.Printf("Status: %d\n\n%s\n", resp.StatusCode, resp.Status)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if x(err, http.StatusInternalServerError) {
		return
	}

	var sparql Sparql
	err = xml.Unmarshal(b, &sparql)
	if x(err, http.StatusInternalServerError) {
		return
	}

	fmt.Printf(`Content-type: text/html; charset=UTF-8

<!DOCTYPE html>
<html lang="nl">
  <head>
    <title>Noordergraf -- %s</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="icon" href="/favicon.ico" type="image/ico">
    <link rel="stylesheet" type="text/css" href="/md.css">
  </head>
  <body class="textsearch">
    <div id="container">
<h1>%s</h1>
gevonden: %d
<table>
`, title, title, len(sparql.Results))

	for _, result := range sparql.Results {
		switch q {
		case "famous":
			var plot, name, sameas string
			for _, binding := range result.Bindings {
				if binding.Name == "plot" {
					plot = binding.URI
				} else if binding.Name == "name" {
					name = binding.Literal
				} else if binding.Name == "sameas" {
					sameas = binding.URI
				}
			}
			plot = strings.Replace(plot, "https://noordergraf.rug.nl/", "", 1)
			fmt.Printf("<tr><td><a href=\"/%s\">%s</a></td><td><a href=\"%s\">%s</a></td></tr>\n", plot, plot, sameas, html.EscapeString(name))
		case "multi":
			var plot, name, sameas string
			for _, binding := range result.Bindings {
				if binding.Name == "plot" {
					plot = binding.URI
				} else if binding.Name == "name" {
					name = binding.Literal
				} else if binding.Name == "sameas" {
					sameas = binding.URI
				}
			}
			plot = strings.Replace(plot, "https://noordergraf.rug.nl/", "", 1)
			sameas = strings.Replace(sameas, "https://noordergraf.rug.nl/", "", 1)
			fmt.Printf("<tr><td><a href=\"/%s\">%s</a></td><td>%s &rarr; <a href=\"/%s\">%s</a></td></tr>\n", plot, plot, html.EscapeString(name), sameas, sameas)
		}
	}

	fmt.Print(`
</table>
</body>
</html>
`)

}

func x(err error, status int, msg ...interface{}) bool {
	if err == nil {
		return false
	}

	var b bytes.Buffer
	_, filename, lineno, ok := runtime.Caller(1)
	if ok {
		b.WriteString(fmt.Sprintf("%v:%v: %v", filename, lineno, err))
	} else {
		b.WriteString(err.Error())
	}
	if len(msg) > 0 {
		b.WriteString(",")
		for _, m := range msg {
			b.WriteString(fmt.Sprintf(" %v", m))
		}
	}
	fmt.Printf(`Status: %d

%s
`, status, b.String())

	return true
}
