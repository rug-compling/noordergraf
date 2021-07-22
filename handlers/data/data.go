package main

/*

TODO:

Voor alle directory's, bijvoorbeeld /pred :

    /pred/
    /pred/index      -> html-bestand met links naar alle bestanden

    /pred/index.ttl
    /pred/index.rdf
    /pred/index.nt   -> alle bestanden samengevoegd

*/

import (
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	HTML = iota
	TURTLE
	TRIPLE
	RDF
)

var (
	format       = HTML
	language     = "en"
	uri          string
	data         string
	prefix       string
	lastModified string
	exts         = map[int]string{
		HTML:   "",
		TURTLE: ".ttl",
		TRIPLE: ".nt",
		RDF:    ".rdf",
	}
	reComment = regexp.MustCompile("(?m:^[ \t]*#.*\n?)")
)

func main() {

	getLanguage()

	// zou niet nodig moeten zijn
	uri = path.Clean(os.Getenv("REQUEST_URI"))
	i := strings.Index(uri, "#")
	if i > 0 {
		uri = uri[:i]
	}

	i = strings.LastIndex(uri, ".")
	if i > 0 {
		switch uri[i:] {
		case ".ttl":
			format = TURTLE
			uri = uri[:i]
		case ".nt":
			format = TRIPLE
			uri = uri[:i]
		case ".rdf":
			format = RDF
			uri = uri[:i]
		}
	}

	filename := "/net/noordergraf/data" + uri + ".ttl"
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print("Status: 404 Not Found\n\n")
		return
	}
	data = string(b)
	b, _ = ioutil.ReadFile("/net/noordergraf/data/prefix.ttl")
	prefix = string(b)
	fi, err := os.Stat(filename)
	if err != nil {
		fmt.Print("Status: 404 Not Found\n\n")
		return
	}
	lastModified = fi.ModTime().UTC().Format(time.RFC1123)

	switch format {
	case HTML:
		doHTML()
	case TURTLE:
		fmt.Print("Content-type: text/turtle; charset=UTF-8\nLast-Modified: " + lastModified + "\n\n")
		fmt.Println(prefix)
		fmt.Print(reComment.ReplaceAllLiteralString(data, ""))
	case TRIPLE:
		fmt.Print("Content-type: application/n-triples; charset=UTF-8\nLast-Modified: " + lastModified + "\n\n")
		fmt.Print(convert("ntriples"))
	case RDF:
		fmt.Print("Content-type: application/rdf+xml; charset=UTF-8\nLast-Modified: " + lastModified + "\n\n")
		fmt.Print(convert("rdfxml-abbrev"))
	}
}

func getLanguage() {
	// HTTP_ACCEPT_LANGUAGE=nl-NL,nl;q=0.9,en;q=0.8

	langs := make(map[string]float64)
	maxval := 0.0
	for _, lang := range strings.Split(os.Getenv("HTTP_ACCEPT_LANGUAGE"), ",") {
		v := 1.0
		aa := strings.Split(lang, ";")
		for _, a := range aa[1:] {
			bb := strings.Split(a, "=")
			if len(bb) == 2 && strings.TrimSpace(bb[0]) == "q" {
				v1, err := strconv.ParseFloat(bb[1], 64)
				if err == nil {
					v = v1
				}
			}
		}
		if v <= maxval {
			continue
		}
		lang1 := strings.ToLower(strings.Split(aa[0], "-")[0])
		langs[lang1] = v
		maxval = v
	}
	if langs["nl"] > langs["en"] {
		language = "nl"
	}
}

func convert(output string) string {
	cmd := exec.Command("rapper", "-i", "turtle", "-o", output, "-I", "https://noordergraf.rug.nl/", "-q", "-")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, prefix)
		io.WriteString(stdin, data)
	}()

	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}

func doHTML() {
	title := html.EscapeString(uri[1:])
	fmt.Printf(`Content-type: text/html; charset=UTF-8
Last-Modified: %s

<html lang="nl">
  <head>
    <title>Noordergraf -- %s</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="icon" href="/favicon.ico" type="image/ico">
    <link rel="stylesheet" type="text/css" href="/md.css">
    <link rel="alternate" href="https://noordergraf.rug.nl%s.ttl" type="text/turtle"/>
    <link rel="alternate" href="https://noordergraf.rug.nl%s.nt"  type="application/n-triples"/>
    <link rel="alternate" href="https://noordergraf.rug.nl%s.rdf" type="application/rdf+xml"/>
  </head>
  <body class="">
    <div id="container">
      <h1>%s</h1>
`, lastModified, title, uri, uri, uri, title)

	fmt.Printf("<pre>\n%s</pre>\n", html.EscapeString(strings.TrimSpace(reComment.ReplaceAllLiteralString(data, ""))))

	fmt.Print(`
    </div>
  </body>
</html>
`)

}
