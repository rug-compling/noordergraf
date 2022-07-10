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
	"strconv"
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
	skips     = make(map[int]bool)
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

	var format, dotDisabled, fdpDisabled, sfdpDisabled string
	switch req.FormValue("f") {
	case "fdp":
		format = "fdp"
		fdpDisabled = " disabled"
	case "sfdp":
		format = "sfdp"
		sfdpDisabled = " disabled"
	default:
		format = "dot"
		dotDisabled = " disabled"
	}

	if req.FormValue("a") != "reset" {
		for _, s := range strings.Split(req.FormValue("s"), "n") {
			if s != "" {
				i, err := strconv.Atoi(s)
				if err == nil {
					skips[i] = true
				}
			}
		}
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
		if !skips[t.id] {
			var e string
			if t.Type == "literal" {
				e = ", shape=box"
			}
			fmt.Fprintf(&buf, "n%d [label=%q%s];\n", t.id, label(t, true), e)
		}
	}

	for _, tri := range data.Triples {
		if tri.Predicate.Value == "http://www.w3.org/1999/02/22-rdf-syntax-ns#type" {
			continue
		}
		if skips[tri.Subject.id] || skips[tri.Object.id] {
			continue
		}
		fmt.Fprintf(&buf, "n%d -> n%d [label=%q];\n", tri.Subject.id, tri.Object.id, label(tri.Predicate, false))
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

	skipobj := make([]string, 0, len(skips))
	for key := range skips {
		skipobj = append(skipobj, fmt.Sprintf(`"n%d":true`, key))
	}

	fmt.Printf(`Content-type: text/html; charset=utf-8

<!DOCTYPE html>
<html>
  <head>
    <title>%s</title>
    <script src="/jquery.js"></script>
    <script>
      var items = {%s};
      function cl(item) {
        var id = $(item).find('title').text();
        if (items[id]) {
          delete items[id];
          item.classList.remove("mark");
        } else {
          items[id] = true;
          item.classList.add("mark");
        }
      }
      function doForm() {
        if (document.forms["myform"]["a"].value != "reset") {
          var s = "";
          for (var key in items) {
            s += key;
          }
          document.forms["myform"]["s"].value = s;
        }
        return true;
      }
    </script>
    <style>
      .mark ellipse,
      .mark polygon {
        stroke: blue;
        stroke-width: 3;
      }
      .scroll {
        padding: 1em 0px;
        margin: 1em 0px;
        overflow-x: auto;
      }
    </style>
  </head>
  <body>
    <form name="myform" action="https://noordergraf.rug.nl/bin/ttl2svg" onsubmit="return doForm()" method="get">
      <input type="hidden" name="t" value="%s">
      <input type="hidden" name="s">
      Lay-out:
      <input type="submit" name="f" value="dot"%s>
      <input type="submit" name="f" value="fdp"%s>
      <input type="submit" name="f" value="sfdp"%s>
      <p>
Klik op tekst in nodes om te selecteren en klik →
      <input type="submit" name="a" value="verbergen">
<p>
      <input type="submit" name="a" value="reset">
    </form>
<div class="scroll">
%s
</div>
    <script>
      $('g.node').on('click', function() { cl(this); });
    </script>
TODO:
<ul>
<li>Bij klik op "verbergen" of "reset" niet terug naar formaat "dot"
<li>Ook dochterelementen verbergen, tenzij die een andere parent hebben die niet verborgen is
</ul>
  </body>
</html>
`, html.EscapeString(tfile), strings.Join(skipobj, ",\n"), tfile, dotDisabled, fdpDisabled, sfdpDisabled, svg)

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
