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
	"regexp"
	"strings"
)

var (
	x      = util.CheckErr
	reDate = regexp.MustCompile(`^[0-9]{4}:[0-9]{2}:[0-9]{2}`)
)

func main() {
	b, err := ioutil.ReadFile(os.Args[1])
	x(err)

	raw, err := exif.SearchAndExtractExif(b)
	x(err)

	entries, _, err := exif.GetFlatExifDataUniversalSearch(raw, nil, true)
	lat := 0.0
	lon := 0.0
	latv := 1.0
	lonv := 1.0
	datetime := ""
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
		case "DateTime":
			if datetime == "" {
				datetime = e.Formatted
			}
		case "DateTimeOriginal":
			if e.Formatted != "" {
				datetime = e.Formatted
			}
		}
	}

	if reDate.MatchString(datetime) {
		datetime = strings.Replace(datetime, ":", "-", 2)
	}
	datetime = strings.Join(strings.Fields(datetime), "T")

	name := filepath.Base(os.Args[1])
	lower := strings.ToLower(name)
	if strings.HasSuffix(lower, ".jpg") {
		name = name[:len(name)-4]
	} else if strings.HasSuffix(lower, ".jpeg") {
		name = name[:len(name)-5]
	}
	fmt.Printf(`tomb:%s.im1 a t:photo.n.01 ;
  # :credit "J. Doe" ;
  :license <https://creativecommons.org/publicdomain/zero/1.0/> ;
  # :license <https://creativecommons.org/licenses/by/4.0/> ;
  # :license <https://creativecommons.org/licenses/by-nc/4.0/> ;
  # :license <https://creativecommons.org/licenses/by-nc-nd/4.0/> ;
  # :license <https://creativecommons.org/licenses/by-nc-sa/4.0/> ;
  # :license <https://creativecommons.org/licenses/by-sa/4.0/> ;
  :file "%s"`, name, filepath.Base(os.Args[1]))
	if datetime != "" {
		fmt.Printf(" ;\n  :date \"%s\"^^xsd:dateTime", datetime)
	}
	if lat != 0 && lon != 0 {
		fmt.Printf(` ;
  geo:lat %.4f ;
  geo:long %.4f ;
  :nd "%+08.4f%+09.4f"^^ll:`,
			lat*latv,
			lon*lonv,
			lat*latv,
			lon*lonv)
	}

	fmt.Print(" .\n\n")

}
