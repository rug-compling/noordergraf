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

	query := `
PREFIX :       <https://noordergraf.rug.nl/ns#>
PREFIX symbol: <https://noordergraf.rug.nl/symbol/>
SELECT ?symbool (count(?symbool) as ?aantal) {
  ?a :symbol / a ?symbool .
}
GROUP BY ?symbool
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
		var name, count string
		for _, binding := range result.Bindings {
			if binding.Name == "symbool" {
				name = binding.URI
				name = name[strings.LastIndex(name, "/")+1:]
			} else if binding.Name == "aantal" {
				count = binding.Literal
			}
		}
		fmt.Printf(`<figure>
  <p><a href="/symbol/%s"><img src="/sym/%s100.png" alt="%s" width="100" height="100"></a></p>
  <figcaption>%s<br>%s</figcaption>
</figure>
`, name, name, name, name, count)
	}

	fmt.Println(`</div>`)

}
