package main

import (
	"github.com/pebbe/util"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

	tags := make(map[string]bool)
	b, err := ioutil.ReadFile("/net/noordergraf/tools/chktags.tags")
	x(err)
	for _, s := range strings.Fields(string(b)) {
		tags[s] = true
	}

	query := `
SELECT DISTINCT (lang(?c) as ?lang) {
  ?a ?b ?c .
  FILTER ( langMatches(lang(?c), "*") )
} ORDER BY ?lang
`

	request := "http://localhost:10035/repositories/noordergraf?infer=false&query=" + url.QueryEscape(query)
	resp, err := http.Get(request)
	x(err)

	if resp.StatusCode > 299 {
		x(fmt.Errorf("Status: %d\n\n%s\n", resp.StatusCode, resp.Status))
	}

	b, err = ioutil.ReadAll(resp.Body)
	x(err)

	var sparql Sparql
	err = xml.Unmarshal(b, &sparql)
	x(err)

	for _, result := range sparql.Results {
		for _, binding := range result.Bindings {
			if binding.Name == "lang" {
				lang := binding.Literal
				if !tags[lang] {
					x(fmt.Errorf("Unknown tag %q", lang))
				}
			}
		}
	}
}
