package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
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

func main() {

	query := `
PREFIX :    <https://noordergraf.rug.nl/ns#>
PREFIX geo: <http://www.w3.org/2003/01/geo/wgs84_pos#>
SELECT ?s ?n ?lat ?lon {
  ?s a :Site .
  ?s :sitename ?n .
  ?s :geo / geo:lat ?lat .
  ?s :geo / geo:long ?lon .
}
`

	request := "http://localhost:10035/repositories/noordergraf?query=" + url.QueryEscape(query)
	resp, err := http.Get(request)
	if x(err, http.StatusInternalServerError) {
		return
	}
	if resp.StatusCode > 299 {
		fmt.Printf("Status: %d\n\n%s\n", resp.StatusCode, resp.Status)
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

	fmt.Print(`Content-type: application/json; encoding=utf-8

[
`)

	p := ""
	for _, result := range sparql.Results {
		var s, n, lat, lon string
		for _, binding := range result.Bindings {
			switch binding.Name {
			case "s":
				s = binding.URI
			case "n":
				n = binding.Literal
			case "lat":
				lat = binding.Literal
			case "lon":
				lon = binding.Literal
			}
		}
		fmt.Printf("%s[%q,%q,%s,%s]", p, s, n, lat, lon)
		p = ",\n"
	}
	fmt.Print("\n]\n")

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
