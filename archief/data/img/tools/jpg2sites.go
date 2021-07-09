package main

import (
	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/pebbe/util"

	"bufio"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	x = util.CheckErr
)

type Site struct {
	x, y, z float64
	lbl     string
}

func main() {

	sites := make([]*Site, 0)
	fp, err := os.Open("../../sites/sites.csv")
	x(err)
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		aa := strings.Split(scanner.Text(), ",")
		if aa[0] == "lat" {
			continue
		}
		lats := strings.Split(aa[0], "/")
		lons := strings.Split(aa[1], "/")
		lat1, err := strconv.ParseFloat(lats[0], 64)
		x(err)
		lat2, err := strconv.ParseFloat(lats[1], 64)
		x(err)
		lat := lat1 / lat2 / 180.0 * math.Pi
		lon1, err := strconv.ParseFloat(lons[0], 64)
		x(err)
		lon2, err := strconv.ParseFloat(lons[1], 64)
		x(err)
		lon := lon1 / lon2 / 180.0 * math.Pi
		sites = append(sites, &Site{
			x:   math.Cos(lon) * math.Cos(lat),
			y:   math.Sin(lon) * math.Cos(lat),
			z:   math.Sin(lat),
			lbl: aa[2][strings.LastIndex(aa[2], "/")+1 : len(aa[2])-2],
		})
	}
	x(scanner.Err())
	fp.Close()

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

		lat = latv * lat / 180.0 * math.Pi
		lon = lonv * lon / 180.0 * math.Pi
		X := math.Cos(lon) * math.Cos(lat)
		Y := math.Sin(lon) * math.Cos(lat)
		Z := math.Sin(lat)

		var nearest *Site
		v := -1.0
		for _, site := range sites {
			v1 := X*site.x + Y*site.y + Z*site.z

			if v1 > v {
				v = v1
				nearest = site
			}
		}
		f := math.Acos(v) / math.Pi * 20000
		if f > .8 {
			continue
		}

		name := filepath.Base(file)
		n := len(name)
		if strings.HasSuffix(name, ".jpg") {
			name = name[:n-4]
		} else if strings.HasSuffix(name, ".jpeg") {
			name = name[:n-5]
		}
		fp, err := os.Create(filepath.Dir(file) + "/../" + name + ".loc")
		x(err)
		fp.WriteString(nearest.lbl + "\n")
		fp.Close()
	}
}
