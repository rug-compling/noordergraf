package main

/*

Dit programma maakt de body van de html-weergave van de lijst met symbolen

Dit programma wordt aangeroepen vanuit /net/noordergraf/data/input

*/

import (
	"github.com/pebbe/util"

	"encoding/xml"
	"fmt"
	"html"
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

	lang := "eng"
	if len(os.Args) > 1 && os.Args[1] == "-n" {
		lang = "nld"
	}

	query := `
PREFIX :       <https://noordergraf.rug.nl/ns#>
PREFIX symbol: <https://noordergraf.rug.nl/symbol/>
SELECT ?symbool (count(?symbool) as ?aantal) ?text {
  ?a :symbol / :symbolType ?symbool .
  ?symbool rdfs:comment ?text .
  FILTER ( langMatches(lang(?text), "` + lang + `") )
}
GROUP BY ?symbool ?text
ORDER BY ?symbool
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

	fmt.Println(`<div class="rows">`)

	for _, result := range sparql.Results {
		var name, title, count string
		for _, binding := range result.Bindings {
			if binding.Name == "symbool" {
				name = binding.URI
				name = name[strings.LastIndex(name, "/")+1:]
			} else if binding.Name == "aantal" {
				count = binding.Literal
			} else if binding.Name == "text" {
				title = html.EscapeString(binding.Literal)
			}
		}
		fmt.Printf(`<figure>
  <p><a href="/symbol/%s"><img src="/picto/%s100.png" alt="%s" title="%s" width="100" height="100"></a></p>
  <figcaption>%s<br>%s</figcaption>
</figure>
`, name, name, title, title, name, count)
	}

	fmt.Println(`</div>`)

}
