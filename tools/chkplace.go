package main

import (
	"github.com/pebbe/util"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	x = util.CheckErr
)

func main() {

	// inference werkt niet binnen een GRAPH blok
	query := `
SELECT ?g ?p ?text {
  ?a a :PlaceClass .
  GRAPH ?g {
    ?s ?p ?a .
  }
  OPTIONAL { ?a :text ?text }
  FILTER NOT EXISTS { ?a :place ?x  }
  FILTER NOT EXISTS { ?a :coordinates ?y  }
} ORDER BY ?g
`

	request := "http://localhost:10035/repositories/noordergraf?infer=true&query=" + url.QueryEscape(query)
	resp, err := http.Get(request)
	x(err)

	if resp.StatusCode > 299 {
		x(fmt.Errorf("Status: %d\n\n%s\n", resp.StatusCode, resp.Status))
	}

	b, err := ioutil.ReadAll(resp.Body)
	x(err)

	var sparql Sparql
	err = xml.Unmarshal(b, &sparql)
	x(err)

	errors := 0
	for _, result := range sparql.Results {
		errors++
		var g, p, text string
		for _, binding := range result.Bindings {
			if binding.Name == "g" {
				g = binding.URI
			} else if binding.Name == "p" {
				p = binding.URI
			} else if binding.Name == "text" {
				text = binding.Literal
			}
		}
		fmt.Fprintln(os.Stderr, "Missing :place and :coordinates in", base(g), base(p), text)
	}
	if errors > 0 {
		os.Exit(1)
	}
}

func base(s string) string {
	i := strings.LastIndexAny(s, "/#")
	if i > 0 {
		return s[i+1:]
	}
	return s
}
