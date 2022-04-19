package main

/*

Dit programma maakt de body van de html-weergave van de lijst met bijbeldelen

Dit programma wordt aangeroepen vanuit /net/noordergraf/data/input

*/

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
	major = map[byte]map[string]string{
		'1': map[string]string{
			"eng": "The Old Testament",
			"nld": "Het Oude Testament",
		},
		'2': map[string]string{
			"eng": "The New Testament",
			"nld": "Het Nieuwe Testament",
		},
	}

	x = util.CheckErr
)

func main() {

	lang := "eng"
	if len(os.Args) > 1 && os.Args[1] == "-n" {
		lang = "nld"
	}

	query := `
PREFIX : <https://noordergraf.rug.nl/ns#>
SELECT ?book ?name ?order (count(?name) as ?cnt)  {
  ?a :book ?book .
  ?book :bibleBookName ?name .
  ?book :order ?order .
  FILTER ( langMatches(lang(?name), "` + lang + `") )
}
GROUP BY ?book ?name ?order
ORDER BY ?order
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

	fmt.Println("<table>")
	if lang == "eng" {
		fmt.Println("<tr><th>number of references<th>book</tr>")
	} else {
		fmt.Println("<tr><th>aantal verwijzingen<th>boek</tr>")
	}
	p := byte('0')
	for _, result := range sparql.Results {
		var book, name, order, cnt string
		for _, binding := range result.Bindings {
			if binding.Name == "book" {
				book = binding.URI
				book = book[strings.LastIndex(book, "/")+1:]
			} else if binding.Name == "name" {
				name = binding.Literal
			} else if binding.Name == "order" {
				order = binding.Literal
			} else if binding.Name == "cnt" {
				cnt = binding.Literal
			}
		}
		if order[0] != p {
			p = order[0]
			fmt.Println(`<tr><th colspan="2">`, major[p][lang], "</tr>")
		}
		fmt.Printf("<tr><td class=\"right\">%s<td><a href=\"/bible/%s\">%s</a></tr>\n", cnt, book, name)
	}
	fmt.Println("</table>")
}
