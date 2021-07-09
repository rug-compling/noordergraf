package main

import (
	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/pebbe/util"

	"bufio"
	"fmt"
	"html"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

var (
	x = util.CheckErr
)

func main() {

	linked := true
	unlinked := true
	title := "graven"

	var files []string

	if os.Args[1] == "-y" {
		unlinked = false
		title = "gelinkte graven"
		files = os.Args[2:]
	} else if os.Args[1] == "-n" {
		linked = false
		title = "vrije graven"
		files = os.Args[2:]
	} else {
		files = os.Args[1:]
	}

	sitemap := make(map[string]string)

	fp, err := os.Open("../../sites/sites.csv")
	x(err)
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		a := strings.Split(scanner.Text(), `>","`)
		if len(a) != 2 {
			continue
		}
		sitemap[a[0][strings.LastIndex(a[0], "/")+1:]] = strings.Trim(a[1], `"`)
	}
	x(scanner.Err())
	fp.Close()

	fmt.Print(`<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://earth.google.com/kml/2.2">
<Document>
<name>` + title + `</name>
`)

	for _, file := range files {

		var loc bool
		var desc string
		if b, err := ioutil.ReadFile("../" + strings.Replace(file, ".jpg", ".loc", 1)); err == nil {
			loc = true
			desc = sitemap[strings.TrimSpace(string(b))]
		}
		if !(loc && linked || !loc && unlinked) {
			continue
		}

		b, err := ioutil.ReadFile(file)
		x(err)

		raw, err := exif.SearchAndExtractExif(b)
		x(err)

		entries, _, err := exif.GetFlatExifDataUniversalSearch(raw, nil, true)
		lat := 0.0
		lon := 0.0
		latv := 1.0
		lonv := 1.0
		for _, e := range entries {
			var f float64
			switch t := e.Value.(type) {
			case []exifcommon.Rational:
				for i, v := range t {
					f += float64(v.Numerator) / float64(v.Denominator) / math.Pow(60.0, float64(i))
				}
			}
			switch e.TagName {
			case "GPSLatitudeRef":
				if e.Formatted == "S" {
					latv = -1.0
				}
			case "GPSLongitudeRef":
				if e.Formatted == "W" {
					lonv = -1.0
				}
			case "GPSLatitude":
				lat = f
			case "GPSLongitude":
				lon = f
			}
		}

		if lat == 0 || lon == 0 {
			continue
		}

		name := filepath.Base(file)
		lower := strings.ToLower(name)
		if strings.HasSuffix(lower, ".jpg") {
			name = name[:len(name)-4]
		} else if strings.HasSuffix(lower, ".jpeg") {
			name = name[:len(name)-5]
		}
		fmt.Printf(`<Placemark>
  <name>%s</name>
  <description>%s</description>
  <Point>
    <coordinates>%.4f,%.4f,0</coordinates>
  </Point>
</Placemark>
`,
			name,
			html.EscapeString(desc),
			lon*lonv,
			lat*latv)
	}
	fmt.Print(`</Document>
</kml>
`)
}
