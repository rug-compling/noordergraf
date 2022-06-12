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
	w *csv.Writer
)

func main() {

	fp, err := os.Create("questions.csv")
	x(err)
	w = csv.NewWriter(fp)

	query := `
SELECT DISTINCT ?item ?photo ?person ?name {
  GRAPH ?item {
    ?person :name / :fullName ?name .
    ?y :file ?photo .
    ?item :todo ?todo .
    VALUES ?todo { todo:ItemType todo:Subjects } .
  }
  FILTER NOT EXISTS { ?z :creator ?person }
}
ORDER BY ?item ?person
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

	var currentItem, currentPhoto string
	var persons []string
	for _, result := range sparql.Results {
		var item, photo, person, name string
		for _, binding := range result.Bindings {
			if binding.Name == "item" {
				item = strings.Replace(binding.URI, "https://noordergraf.rug.nl/item/", "", 1)
			} else if binding.Name == "photo" {
				photo = strings.Replace(binding.URI, "https://noordergraf.rug.nl/img/", "", 1)
			} else if binding.Name == "person" {
				person = binding.URI
			} else if binding.Name == "name" {
				name = strings.TrimSpace(binding.Literal)
			}
		}
		if item != currentItem {
			if currentItem != "" {
				output(currentItem, currentPhoto, persons)
			}
			currentItem = item
			currentPhoto = photo
			persons = make([]string, 0)
		}
		if i := strings.Index(person, "#"); i > 0 {
			persons = append(persons, person[i+1:]+"#"+name)
		}
	}
	output(currentItem, currentPhoto, persons)
	w.Flush()
	fp.Close()
}

func output(item string, photo string, persons []string) {
	plist := strings.Join(persons, "|")
	x(w.Write([]string{item, photo, plist}))
}
