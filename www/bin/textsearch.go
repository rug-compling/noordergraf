package main

import (
	"github.com/rug-compling/noordergraf/go/nlsoundex"

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
	if q == "" {
		fmt.Print(`Status: 400

Missing query
`)
		return
	}
	qq := q

	var template string
	switch req.FormValue("t") {
	case "fullname":
		template = `
PREFIX : <https://noordergraf.rug.nl/ns#>
SELECT ?s ?o {
  ?s :name ?n .
  ?n fti:match ( %q "fullname" ) .
  ?n :fullName ?o .
}
ORDER BY ?s
`
		qq = nlsoundex.Soundex(q)
	default:
		template = `
PREFIX : <https://noordergraf.rug.nl/ns#>
SELECT DISTINCT ?s ?o {
  GRAPH ?s {
    (?s ?o) fti:match ( %q "all" ) .
    ?s :text ?o .
  }
}
ORDER BY ?s ?o
`
	}

	query := fmt.Sprintf(template, qq)

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
    <title>Noordergraf -- zoeken</title>
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
`, html.EscapeString(q), len(sparql.Results))

	for _, result := range sparql.Results {
		var uri, obj string
		for _, binding := range result.Bindings {
			if binding.Name == "s" {
				uri = binding.URI
			} else if binding.Name == "o" {
				obj = binding.Literal
			}
		}
		uri = strings.Replace(uri, "https://noordergraf.rug.nl/", "", 1)
		fmt.Printf("<tr><td><a href=\"/%s\">%s</a></td><td>%s</td></tr>\n", uri, uri, strings.Replace(html.EscapeString(obj), "\n", "<br>", -1))
	}

	fmt.Print(`</dl>
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
