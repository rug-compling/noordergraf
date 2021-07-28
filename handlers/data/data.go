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
	pretab       = make(map[string]string)
	lastModified string
	exts         = map[int]string{
		HTML:   "",
		TURTLE: ".ttl",
		TRIPLE: ".nt",
		RDF:    ".rdf",
	}
	reTokens = regexp.MustCompile(strings.Join([]string{
		"[ \t\n\r\f]+",
		"(&(#34|quot);){3}(.|\n)*?(&(#34|quot);){3}[^ \t\n\r\f]*",
		"&(#34|quot);(.*?)&(#34|quot);[^ \t\n\r\f]*",
		"&lt;http.*?&gt;",
		"[A-Za-z0-9.]*:[A-Za-z0-9.]+",
		".",
	}, "|"))
	reComment = regexp.MustCompile("(?m:^[ \t]*#.*\n?)")
	reND      = regexp.MustCompile(`:nd +&(#34|quot);([-+][.0-9]+)([-+][.0-9]+)&(#34|quot);\^\^ll:`)
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
	prefix = trimPrefix(string(b), data)
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
	body := html.EscapeString(strings.TrimSpace(reComment.ReplaceAllLiteralString(data, "")))

	if uri == "/ns" || strings.HasPrefix(uri, "/place/") || strings.HasPrefix(uri, "/site/") || strings.HasPrefix(uri, "/tomb/") {
		body = reND.ReplaceAllStringFunc(body, func(s string) string {
			m := reND.FindStringSubmatch(s)
			return fmt.Sprintf(`:nd &%s;<a href="geo:%s,%s">%s%s</a>&%s;^^ll:`, m[1], getLL(m[2], 2), getLL(m[3], 3), m[2], m[3], m[4])
		})

		inPrefix := false
		tokens := reTokens.FindAllString(body, -1)
		for i, token := range tokens {
			if token == "@" {
				inPrefix = true
				continue
			}
			if strings.Contains(token, "\n") {
				inPrefix = false
				continue
			}
			if strings.HasPrefix(token, "&lt;http") && strings.HasSuffix(token, "&gt;") {
				if !inPrefix {
					s := token[4 : len(token)-4]
					tokens[i] = `&lt;<a href="` + s + `">` + s + `</a>&gt;`
				}
				continue
			}
			if i > 0 && !strings.HasSuffix(tokens[i-1], "\n") {
				a := strings.SplitN(token, ":", 2)
				if len(a) == 2 && a[1] != "" {
					if s, ok := pretab[a[0]]; ok {
						tokens[i] = `<a href="` + s + a[1] + `">` + token + `</a>`
						continue
					}
				}
			}
		}
		body = strings.Join(tokens, "")
	}

	if uri == "/ns" {
		lines := strings.Split(body, "\n")
		for i, line := range lines {
			if !strings.HasPrefix(line, ":") {
				continue
			}
			a := strings.SplitN(line, " ", 2)
			a[0] = fmt.Sprintf("<span id=%q class=\"hash\">%s</span>", a[0][1:], a[0])
			lines[i] = strings.Join(a, " ")
		}
		body = strings.Join(lines, "\n")
	}

	title := html.EscapeString(uri[1:])
	if title == "ns" {
		title = "ns#"
	}
	fmt.Printf(`Content-type: text/html; charset=UTF-8
Last-Modified: %s

<html lang="nl">
  <head>
    <title>Noordergraf -- %s</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="icon" href="/favicon.ico" type="image/ico">
    <link rel="stylesheet" type="text/css" href="/data.css">
    <link rel="alternate" href="https://noordergraf.rug.nl%s.ttl" type="text/turtle"/>
    <link rel="alternate" href="https://noordergraf.rug.nl%s.nt"  type="application/n-triples"/>
    <link rel="alternate" href="https://noordergraf.rug.nl%s.rdf" type="application/rdf+xml"/>
  </head>
  <body class="">
    <div id="container">
      <h1>%s</h1>
`, lastModified, title, uri, uri, uri, title)

	fmt.Printf("<pre>\n%s\n\n%s\n</pre>\n", html.EscapeString(strings.TrimSpace(prefix)), body)

	fmt.Print(`
    </div>
  </body>
</html>
`)

}

func getLL(s string, n int) string {
	sign := 1.0
	if s[0] == '-' {
		sign = -1.0
	}

	var digits string
	i := strings.Index(s, ".")
	if i > 0 {
		digits = s[1:i]
	} else {
		digits = s[1:]
	}
	ln := len(digits)

	v, _ := strconv.Atoi(digits[:n])
	value := float64(v)
	if ln >= n+2 {
		v, _ := strconv.Atoi(digits[n : n+2])
		value += float64(v) / 60.0
	}
	if ln >= n+4 {
		v, _ := strconv.Atoi(digits[n+2 : n+4])
		value += float64(v) / 3600.0
	}

	if i > 0 {
		v, _ := strconv.ParseFloat(s[i:], 64)
		if ln == n {
			value += v
		} else if ln == n+2 {
			value += v / 60.0
		} else if ln == n+4 {
			value += v / 3600.0
		}
	}

	return fmt.Sprintf("%.4f", value*sign)
}

func trimPrefix(prefix, data string) string {
	lines := make([]string, 0)
	for _, line := range strings.SplitAfter(prefix, "\n") {
		a := strings.Fields(line)
		if len(a) > 2 && strings.Contains(data, a[1]) {
			lines = append(lines, line)
			key := a[1][:len(a[1])-1]
			if key != "rdf" && key != "rdfs" && key != "xsd" {
				pretab[key] = a[2][1 : len(a[2])-1]
			}
		}
	}
	return strings.Join(lines, "")
}
