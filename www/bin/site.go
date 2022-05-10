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
	"os"
	"regexp"
	"runtime"
	"strconv"
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

var (
	reSite = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	lang   = 1
)

func main() {
	getLanguage()

	req, err := cgi.Request()
	if x(err, http.StatusInternalServerError) {
		return
	}
	if x(req.ParseForm(), http.StatusInternalServerError) {
		return
	}
	site := req.FormValue("q")
	if site == "" {
		fmt.Print(`Status: 400

Missing query
`)
		return
	}

	if !reSite.MatchString(site) {
		fmt.Print(`Status: 400

Invalid query
`)
		return
	}

	query := `
PREFIX :       <https://noordergraf.rug.nl/ns#>
PREFIX site:   <https://noordergraf.rug.nl/site/>
SELECT DISTINCT ?tomb ?name {
  ?tomb :site site:` + site + ` .
  ?tomb :subject ?p .
  ?p a :Person .
  ?p :name / :fullName ?name .
}
ORDER BY ?name ?tomb
`

	request := "http://localhost:10035/repositories/noordergraf?infer=true&query=" + url.QueryEscape(query)
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

	langtag := []string{"nl", "en"}[lang]
	title := []string{"Graven op", "Graves at"}[lang]
	found := []string{"gevonden", "found"}[lang]

	fmt.Printf(`Content-type: text/html; charset=UTF-8

<!DOCTYPE html>
<html lang="%s">
  <head>
    <title>Noordergraf -- %s op %s</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="icon" href="/favicon.ico" type="image/ico">
    <link rel="stylesheet" type="text/css" href="/md.css">
  </head>
  <body class="textsearch">
    <div id="container">
<h1>%s %s</h1>
%s: %d
<table>
`, langtag, title, site, title, site, found, len(sparql.Results))

	for _, result := range sparql.Results {
		var uri, name string
		for _, binding := range result.Bindings {
			if binding.Name == "tomb" {
				uri = binding.URI
			} else if binding.Name == "name" {
				name = binding.Literal
			}
		}
		uri = strings.Replace(uri, "https://noordergraf.rug.nl/", "", 1)
		fmt.Printf("<tr><td><a href=\"/%s\">%s</a></td><td>%s</td></tr>\n", uri, uri, strings.Replace(html.EscapeString(name), "\n", "<br>", -1))
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

func getLanguage() {
	// HTTP_ACCEPT_LANGUAGE=nl-NL,nl;q=0.9,en;q=0.8

	langs := make(map[string]float64)
	maxval := 0.0
	for _, lang := range strings.Split(os.Getenv("HTTP_ACCEPT_LANGUAGE"), ",") {
		v := 1.0
		aa := strings.Split(lang, ";")
		for _, a := range aa[1:] {
			bb := strings.Split(a, "=")
			if len(bb) == 2 && strings.TrimSpace(bb[0]) == "q" {
				v1, err := strconv.ParseFloat(bb[1], 64)
				if err == nil {
					v = v1
				}
			}
		}
		if v <= maxval {
			continue
		}
		lang1 := strings.ToLower(strings.Split(aa[0], "-")[0])
		langs[lang1] = v
		maxval = v
	}
	if langs["nl"] > langs["en"] {
		lang = 0
	}
}
