package main

import (
	"github.com/pebbe/util"
	"github.com/rug-compling/noordergraf/go/nlsoundex"
	"github.com/umahmood/soundex"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var (
	lines = make([]string, 0)
	re    = regexp.MustCompile(`[a-zA-Z][a-zA-Z]+`)
	x     = util.CheckErr
)

func main() {
	standard := len(os.Args) > 1 && os.Args[1] == "-s"
	trnsfrm := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		var s string
		if standard {
			s, _, _ = transform.String(trnsfrm, line)
			s = re.ReplaceAllStringFunc(s, func(s1 string) string {
				return soundex.Code(s1)
			})
		} else {
			s = nlsoundex.Soundex(line)
		}
		lines = append(lines, s+"\t"+line)
	}
	x(scanner.Err())
	sort.Strings(lines)
	p := "A0000"
	for _, line := range lines {
		if e := strings.Split(line, "\t")[0]; e != p {
			p = e
			fmt.Println()
		}
		fmt.Println(line)
	}
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}
