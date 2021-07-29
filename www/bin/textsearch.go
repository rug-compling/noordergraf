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
	URIs []string `xml:"results>result>binding>uri"`
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

	query := fmt.Sprintf(`
SELECT DISTINCT ?graph {
  GRAPH ?graph {
    ?s fti:match %q .
  }
}
ORDER BY ?graph
`, q)

	request := "http://localhost:10035/repositories/noordergraf?limit=1000&query=" + url.QueryEscape(query)
	resp, err := http.Get(request)
	if x(err, http.StatusInternalServerError) {
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

<html lang="nl">
  <head>
    <title>Noordergraf -- zoeken</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="icon" href="/favicon.ico" type="image/ico">
    <link rel="stylesheet" type="text/css" href="/md.css">
  </head>
  <body class="">
    <div id="container">
<h1>%s</h1>
Gevonden: %d
<ol>
`, html.EscapeString(q), len(sparql.URIs))

	for _, uri := range sparql.URIs {
		uri = strings.Replace(uri, "https://noordergraf.rug.nl/", "", 1)
		fmt.Printf("<li><a href=\"/%s\">%s</a>\n", uri, uri)
	}

	fmt.Print(`</ol>
</div>
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
