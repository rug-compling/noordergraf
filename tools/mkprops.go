package main

import (
	"github.com/pebbe/util"

	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

	classes := make(map[string]map[string][]string)

	query := `
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
SELECT ?class {
  ?class a rdfs:Class .
}
ORDER BY ?class
`

	request := "http://localhost:10035/repositories/noordergraf?infer=false&query=" + url.QueryEscape(query)
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

	for _, result := range sparql.Results {
		var class string
		for _, binding := range result.Bindings {
			if binding.Name == "class" {
				class = binding.URI
				classes[class] = make(map[string][]string)
				classes[class]["sub"] = make([]string, 0)
				classes[class]["prop"] = make([]string, 0)
			}
		}
	}

	for class := range classes {
		query := `
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
SELECT ?sub {
  ?sub rdfs:subClassOf <` + class + `> .
} ORDER BY ?sub
`

		request := "http://localhost:10035/repositories/noordergraf?infer=false&query=" + url.QueryEscape(query)
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

		for _, result := range sparql.Results {
			var sub string
			for _, binding := range result.Bindings {
				if binding.Name == "sub" {
					sub = binding.URI
					classes[class]["sub"] = append(classes[class]["sub"], sub)
				}
			}
		}
	}

	for class := range classes {
		query := `
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
SELECT ?prop {
  ?prop rdfs:domain / ^rdfs:subClassOf* <` + class + `> .
} ORDER BY ?prop
`

		request := "http://localhost:10035/repositories/noordergraf?infer=false&query=" + url.QueryEscape(query)
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

		for _, result := range sparql.Results {
			var prop string
			for _, binding := range result.Bindings {
				if binding.Name == "prop" {
					prop = binding.URI
					classes[class]["prop"] = append(classes[class]["prop"], prop)
				}
			}
		}
	}

	classlist := make([]string, 0)
	for class := range classes {
		classlist = append(classlist, class)
	}
	for _, class := range classlist {
		if len(classes[class]["sub"]) == 0 && len(classes[class]["prop"]) == 0 {
			delete(classes, class)
		}
	}

	b, err = json.MarshalIndent(classes, "", "  ")
	x(err)
	fmt.Println(string(b))

}
