package main

import (
	fhtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	mdhtml "github.com/yuin/goldmark/renderer/html"

	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	linenumbers = false

	reLi   = regexp.MustCompile(`<li>(\.[^ ]*|\[[ xX]?\]) `)
	reLang = regexp.MustCompile(`{{(.*?)\|(.*?)}}`)

	li = map[string]string{
		`<li>[] `:    `<li class="todo">`,
		`<li>[ ] `:   `<li class="todo">`,
		`<li>[x] `:   `<li class="done">`,
		`<li>[X] `:   `<li class="done">`,
		`<li>. `:     `<li>`,
		`<li>.] `:    `<li class="vak">`,
		`<li>.P `:    `<li class="P">`,
		`<li>.V `:    `<li class="V">`,
		`<li>.v `:    `<li class="v">`,
		`<li>./ `:    `<li class="half">`,
		`<li>.X `:    `<li class="X">`,
		`<li>.&gt; `: `<li class="pijl">`,
		`<li>.# `:    `<li class="ast">`,
		`<li>.+ `:    `<li class="ster">`,
		`<li>.B `:    `<li class="boek">`,
		`<li>.? `:    `<li class="vraag">`,
		`<li>.! `:    `<li class="imp">`,
	}

	header = `Content-type: text/html; charset=utf-8
Cache-Control: public, max-age=86400

<!DOCTYPE html>
<html lang="nl">
  <head>
    <title>TITLE</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="icon" href="/favicon.ico" type="image/ico">
    <link rel="stylesheet" type="text/css" href="/STYLE">
<!--HEAD-->  </head>
  <body class="CLASS">
    <div id="container">
`

	footer = `
    </div>
  </body>
</html>
`
)

func main() {

	//	parts := bytes.Split(

	file := os.Getenv("PATH_TRANSLATED")

	if !strings.HasPrefix(file, "/net/noordergraf/www/") {
		fmt.Print("Status: 404\n\n")
		return
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print("Status: 500\nContent-type: text/plain\n\n", err, "\n")
		return
	}

	b = selectLanguage(b)

	parts := bytes.SplitN(b, []byte("////"), 2)
	if len(parts) == 2 {
		b = parts[1]
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough,
			extension.Linkify,
			extension.DefinitionList,
			extension.NewTypographer(
				// alleen: -- ---
				extension.WithTypographicSubstitutions(extension.TypographicSubstitutions{
					extension.LeftSingleQuote:  nil,
					extension.RightSingleQuote: nil,
					extension.LeftDoubleQuote:  nil,
					extension.RightDoubleQuote: nil,
					extension.Ellipsis:         nil,
					extension.LeftAngleQuote:   nil,
					extension.RightAngleQuote:  nil,
					extension.Apostrophe:       nil,
				}),
			),
			highlighting.NewHighlighting(
				highlighting.WithStyle("manni"),
				highlighting.WithFormatOptions(
					fhtml.WithLineNumbers(linenumbers),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			mdhtml.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	err = md.Convert(b, &buf)
	if err != nil {
		fmt.Print("Status: 500\nContent-type: text/plain\n\n", err, "\n")
		return
	}

	bs := strings.Replace(buf.String(), "<ul>\n<li>.", "<ul class=\"empty\">\n<li>.", -1)
	bs = strings.Replace(bs, "<ul>\n<li>[", "<ul class=\"empty\">\n<li>[", -1)
	bs = reLi.ReplaceAllStringFunc(bs, func(s string) string {
		if s2, ok := li[s]; ok {
			return s2
		}
		return s
	})

	bs = strings.Replace(bs, "<table>", `<div class="overflow mytable"><table>`, -1)
	bs = strings.Replace(bs, "</table>", `</table></div>`, -1)

	title := os.Getenv("SCRIPT_URI")
	head := ""
	style := "md.css"
	class := ""
	if len(parts) > 1 {
		for _, line := range strings.Split(string(parts[0]), "\n") {
			aa := strings.SplitN(line, ":", 2)
			if len(aa) == 2 {
				value := strings.TrimSpace(aa[1])
				if value != "" {
					switch strings.TrimSpace(aa[0]) {
					case "title":
						title = value
					case "style":
						style = value
					case "class":
						class = value
					case "head":
						head += "    " + value + "\n"
					}
				}
			}
		}
	}

	fmt.Print(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(header, "TITLE", title, 1),
					"STYLE", style, 1),
				"CLASS", class, 1),
			"<!--HEAD-->", head, 1),
		bs,
		footer)
}

func selectLanguage(b []byte) []byte {
	switch getLanguage() {
	case "nl":
		b = reLang.ReplaceAll(b, []byte("$1"))
	default:
		b = reLang.ReplaceAll(b, []byte("$2"))
	}
	return b
}

func getLanguage() string {
	// HTTP_ACCEPT_LANGUAGE=nl-NL,nl;q=0.9,en;q=0.8

	language := "en"
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
	return language
}
