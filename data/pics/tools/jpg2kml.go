package main

import (
	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/pebbe/util"

	"fmt"
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

	fmt.Print(`<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://earth.google.com/kml/2.2">
<Document>
<name>tombreader</name>
`)

	for _, file := range os.Args[1:] {

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
  <Point>
    <coordinates>%.4f,%.4f,0</coordinates>
  </Point>
</Placemark>
`,
			name,
			lon*lonv,
			lat*latv)
	}
	fmt.Print(`</Document>
</kml>
`)
}
