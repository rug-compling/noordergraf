package main

import (
	"github.com/pebbe/util"

	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
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
	x      = util.CheckErr
	w      *csv.Writer
	reTodo = regexp.MustCompile(`^[\t\r\n ]*\b[tT][oO][dD][oO]\b[:\t\r\n ]*`)
)

func main() {

	fp, err := os.Create("questions.csv")
	x(err)
	w = csv.NewWriter(fp)

	query := `
SELECT DISTINCT ?item ?photo ?text {
  GRAPH ?item {
    ?y :file ?photo .
    ?item :todo todo:Text .
    ?item :text ?text .
  }
}
ORDER BY ?item
`

	request := "http://localhost:10035/repositories/noordergraf?limit=10000&query=" + url.QueryEscape(query)
	resp, err := http.Get(request)
	x(err)
	if resp.StatusCode > 299 {
		x(fmt.Errorf("Status: %d\n\n%s\n", resp.StatusCode, resp.Status))
	}

	b, err := ioutil.ReadAll(resp.Body)
	x(err)

	var sparql Sparql
	x(xml.Unmarshal(b, &sparql))

	for _, result := range sparql.Results {
		var item, photo, text string
		for _, binding := range result.Bindings {
			if binding.Name == "item" {
				item = strings.Replace(binding.URI, "https://noordergraf.rug.nl/item/", "", 1)
			} else if binding.Name == "photo" {
				photo = strings.Replace(binding.URI, "https://noordergraf.rug.nl/img/", "", 1)
			} else if binding.Name == "text" {
				text = strings.TrimSpace(reTodo.ReplaceAllLiteralString(binding.Literal, ""))
			}
		}
		x(w.Write([]string{item, photo, text}))
	}
	w.Flush()
	fp.Close()
}
