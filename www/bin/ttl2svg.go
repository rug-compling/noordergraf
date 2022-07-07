package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Data struct {
	Triples []TripleT `json:"triples"`
}

type TripleT struct {
	Subject   *ThingT `json:"subject"`
	Predicate *ThingT `json:"predicate"`
	Object    *ThingT `json:"object"`
}

type ThingT struct {
	Value    string `json:"value"`
	Datatype string `json:"datatype"`
	Type     string `json:"type"`
	id       int
}

var (
	nid       = 0
	thingMap  = make(map[string]int)
	thingList = make([]*ThingT, 0)
	types     = make(map[int]string)
	short     = make([][2]string, 0)
)

func main() {

	req, err := cgi.Request()
	if x(err, http.StatusInternalServerError) {
		return
	}
	if x(req.ParseForm(), http.StatusInternalServerError) {
		return
	}
	tfile := req.FormValue("t")
	if tfile == "" {
		fmt.Print(`Status: 400

Missing query
`)
		return
	}

	format := "dot"
	switch req.FormValue("f") {
	case "fdp":
		format = "fdp"
	case "sfdp":
		format = "sfdp"
	}

	fp, err := os.Open("/net/noordergraf/data/prefix.ttl")
	if x(err, http.StatusInternalServerError) {
		return
	}
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "@prefix") {
			continue
		}
		aa := strings.Fields(line)
		short = append(short, [2]string{strings.Replace(aa[2][1:len(aa[2])-1], "https://noordergraf.rug.nl/", "", 1), aa[1]})
	}
	if x(scanner.Err(), http.StatusInternalServerError) {
		return
	}
	fp.Close()

	filename := "/net/noordergraf/data/" + strings.Replace(tfile, ":", "/", 1) + ".ttl"

	var buf bytes.Buffer
	b, err := ioutil.ReadFile("/net/noordergraf/data/prefix.ttl")
	if x(err, 404) {
		return
	}
	buf.Write(b)
	b, err = ioutil.ReadFile(filename)
	if x(err, 404) {
		return
	}
	buf.Write(b)
	cmd := exec.Command("rapper", strings.Fields("-i turtle -o json-triples -I https://noordergraf.rug.nl/ -q -")...)
	cmd.Stdin = &buf
	b, err = cmd.Output()
	if x(err, http.StatusInternalServerError) {
		return
	}

	var data Data
	if x(json.Unmarshal(b, &data), http.StatusInternalServerError) {
		return
	}

	for _, tri := range data.Triples {
		putThing(tri.Subject)
		if tri.Predicate.Value == "http://www.w3.org/1999/02/22-rdf-syntax-ns#type" {
			types[tri.Subject.id] = label(tri.Object, false)
		} else {
			putThing(tri.Object)
		}
	}

	buf.Reset()

	fmt.Fprintln(&buf, "digraph gr {")

	for _, t := range thingList {
		var e string
		if t.Type == "literal" {
			e = ", shape=box"
		}
		fmt.Fprintf(&buf, "n%02d [label=%q%s];\n", t.id, label(t, true), e)
	}

	for _, tri := range data.Triples {
		if tri.Predicate.Value == "http://www.w3.org/1999/02/22-rdf-syntax-ns#type" {
			continue
		}
		fmt.Fprintf(&buf, "n%02d -> n%02d [label=%q];\n", tri.Subject.id, tri.Object.id, label(tri.Predicate, false))
	}

	fmt.Fprintln(&buf, "}")

	cmd = exec.Command(format, "-Tsvg")
	cmd.Stdin = &buf
	b, err = cmd.Output()
	if x(err, http.StatusInternalServerError) {
		return
	}

	svg := string(b)
	svg = svg[strings.Index(svg, "<svg"):]

	fmt.Printf(`Content-type: text/html; charset=utf-8

<!DOCTYPE html>
<html>
  <head>
    <title>%s</title>
  </head>
  <body>
%s
  </body>
</html>
`, html.EscapeString(tfile), svg)

}

func putThing(t *ThingT) {
	if t.Type == "literal" {
		t.id = nid
		nid++
		thingList = append(thingList, t)
		return
	}
	id, ok := thingMap[t.Value]
	if ok {
		t.id = id
		return
	}
	t.id = nid
	nid++
	thingList = append(thingList, t)
	thingMap[t.Value] = t.id
	return
}

func label(t *ThingT, withType bool) string {
	if t.Type == "literal" {
		return t.Value
	}
	var lbl string
	if t.Type == "bnode" {
		lbl = "--"
	} else {
		lbl = t.Value
		for _, s := range short {
			if strings.HasPrefix(lbl, s[0]) {
				lbl = strings.Replace(lbl, s[0], s[1], 1)
				break
			}
		}
	}
	if withType {
		if tp, ok := types[t.id]; ok {
			lbl += "\n" + tp
		}
	}
	return lbl
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
