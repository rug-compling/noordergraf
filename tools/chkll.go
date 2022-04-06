package main

/*

Dit programma controleert dat de predicaten :geo en :nd altijd samen worden gebruikt
en dat de waardes met elkaar overeen komen.

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
	x = util.CheckErr
)

func main() {

	pattern := "%+08.4f%+09.4f"

	// match tussen :geo en :nd

	query := `
PREFIX :    <https://noordergraf.rug.nl/ns#>
PREFIX geo: <http://www.w3.org/2003/01/geo/wgs84_pos#>
SELECT ?g ?ll ?lat ?lon {
  GRAPH ?g {
    ?s :geo ?geo .
    ?geo geo:lat ?lat .
    ?geo geo:long ?lon .
    ?s :nd ?ll .
  }
} ORDER BY ?g
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

	errors := 0
	for _, result := range sparql.Results {
		var graph, ll, lat, lon string
		for _, binding := range result.Bindings {
			if binding.Name == "g" {
				graph = binding.URI
			} else if binding.Name == "ll" {
				ll = binding.Literal
			} else if binding.Name == "lat" {
				lat = binding.Literal
			} else if binding.Name == "lon" {
				lon = binding.Literal
			}
		}
		i := strings.IndexAny(ll[1:], "-+") + 1
		lat1, _ := strconv.ParseFloat(ll[:i], 64)
		lat2, _ := strconv.ParseFloat(lat, 64)
		lon1, _ := strconv.ParseFloat(ll[i:], 64)
		lon2, _ := strconv.ParseFloat(lon, 64)
		ll1 := fmt.Sprintf(pattern, lat1, lon1)
		ll2 := fmt.Sprintf(pattern, lat2, lon2)
		if ll1 != ll2 {
			errors++
			fmt.Fprintln(os.Stderr, "ERROR", graph[27:], ll1, lat, lon)
		}
	}

	// :geo zonder :nd

	query = `
PREFIX :    <https://noordergraf.rug.nl/ns#>
PREFIX geo: <http://www.w3.org/2003/01/geo/wgs84_pos#>
SELECT ?g ?lat ?lon {
  GRAPH ?g {
    ?s :geo ?geo .
    ?geo geo:lat ?lat .
    ?geo geo:long ?lon .
    FILTER NOT EXISTS { ?s :nd ?ll }
  }
} ORDER BY ?g
`

	request = "http://localhost:10035/repositories/noordergraf?infer=false&query=" + url.QueryEscape(query)
	resp, err = http.Get(request)
	x(err)

	if resp.StatusCode > 299 {
		x(fmt.Errorf("Status: %d\n\n%s\n", resp.StatusCode, resp.Status))
	}

	b, err = ioutil.ReadAll(resp.Body)
	x(err)

	sparql.Results = sparql.Results[0:0]
	err = xml.Unmarshal(b, &sparql)
	x(err)

	for _, result := range sparql.Results {
		var graph, lat, lon string
		for _, binding := range result.Bindings {
			if binding.Name == "g" {
				graph = binding.URI
			} else if binding.Name == "lat" {
				lat = binding.Literal
			} else if binding.Name == "lon" {
				lon = binding.Literal
			}
		}
		errors++
		fmt.Fprintln(os.Stderr, "MISSING :nd", graph[27:], lat, lon)
	}

	// :nd zonder :geo

	query = `
PREFIX : <https://noordergraf.rug.nl/ns#>
SELECT ?g ?ll {
  GRAPH ?g {
    ?s :nd ?ll .
    FILTER NOT EXISTS { ?s :geo ?geo }
  }
} ORDER BY ?g
`

	request = "http://localhost:10035/repositories/noordergraf?infer=false&query=" + url.QueryEscape(query)
	resp, err = http.Get(request)
	x(err)

	if resp.StatusCode > 299 {
		x(fmt.Errorf("Status: %d\n\n%s\n", resp.StatusCode, resp.Status))
	}

	b, err = ioutil.ReadAll(resp.Body)
	x(err)

	sparql.Results = sparql.Results[0:0]
	err = xml.Unmarshal(b, &sparql)
	x(err)

	for _, result := range sparql.Results {
		var graph, nd string
		for _, binding := range result.Bindings {
			if binding.Name == "g" {
				graph = binding.URI
			} else if binding.Name == "nd" {
				nd = binding.Literal
			}
		}
		errors++
		fmt.Fprintln(os.Stderr, "MISSING :geo", graph[27:], nd)
	}

	if errors > 0 {
		os.Exit(1)
	}
}
