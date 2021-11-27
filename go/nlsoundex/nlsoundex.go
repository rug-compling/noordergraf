package nlsoundex

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"regexp"
	"sort"
	"strings"
	"unicode"
)

type sub struct {
	r *regexp.Regexp
	s string
}

var (
	trnsfrm transform.Transformer

	reWord = regexp.MustCompile(`IJ[A-Z]+|I[A-IK-Z][A-Z]*|[A-HJ-Z][A-Z]+`)

	specials = []sub{
		// verbindings-s
		sub{
			r: regexp.MustCompile(`(..)S([^AEIOUYPT])`),
			s: `${1}${2}`,
		},
		// verbindings-n en -r
		// slot-n, -r, -ns en -rs
		sub{
			r: regexp.MustCompile(`([^AEIOUYJ]E)[NR](S?$|([^AEIOUY]))`),
			s: `${1}${3}`,
		},
	}

	reCode *regexp.Regexp
	reDup  = regexp.MustCompile(`11+|22+|33+|44+|55+|66+|77+|88+|99+`)

	cc = map[string]string{

		"B": "1",
		"P": "1",
		"F": "1",
		"V": "1",
		"W": "1",

		"G": "2",
		"K": "2",
		"Q": "2",

		"D": "3",
		"T": "3",

		"L": "4",

		"M": "5",
		"N": "5",

		"R": "6",

		"S": "7",
		"X": "7",
		"Z": "7",
	}
)

func init() {
	trnsfrm = transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	keys := make([]string, 0)
	for key := range cc {
		if len(key) > 1 {
			keys = append(keys, key)
		}
	}
	// omkeerd, zodat CH voor C komt
	sort.Slice(keys, func(i, j int) bool {
		return keys[j] < keys[i]
	})

	keys = append(keys, ".")

	reCode = regexp.MustCompile(strings.Join(keys, "|"))
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func Soundex(txt string) string {

	// accenten weg
	txt, _, _ = transform.String(trnsfrm, txt) // TODO: error

	txt = strings.ToUpper(txt)

	return reWord.ReplaceAllStringFunc(txt, func(s string) string {

		// klankvarianten Nederlandse spelling
		// veel van deze omzettingen zijn alleen van belang voor de beginletter

		s = strings.ReplaceAll(s, "EI", "Y")
		s = strings.ReplaceAll(s, "EY", "Y")
		s = strings.ReplaceAll(s, "IJ", "Y")
		if strings.HasPrefix(s, "CHR") {
			s = "K" + s[2:]
		}
		s = strings.ReplaceAll(s, "SCH", "SS") // g-klank genegeerd
		s = strings.ReplaceAll(s, "CH", "GG")
		s = strings.ReplaceAll(s, "CE", "SE")
		s = strings.ReplaceAll(s, "CI", "SI")
		s = strings.ReplaceAll(s, "CY", "SY")
		s = strings.ReplaceAll(s, "C", "K")

		s = strings.ReplaceAll(s, "X", "SS") // k-klank genegeerd
		s = strings.ReplaceAll(s, "Z", "S")

		s = strings.ReplaceAll(s, "V", "F")

		s = strings.ReplaceAll(s, "AUT", "ALT")
		s = strings.ReplaceAll(s, "OUD", "OLD") // meer ?

		s = strings.ReplaceAll(s, "OU", "AU")

		for _, sp := range specials {
			s = sp.r.ReplaceAllString(s, sp.s)
		}

		se := reCode.ReplaceAllStringFunc(s, func(s1 string) string {
			if s2, ok := cc[s1]; ok {
				return s2
			}
			return "0"
		})

		se = reDup.ReplaceAllStringFunc(se, func(s1 string) string {
			return s1[:1]
		})

		se = se[:1] + strings.ReplaceAll(se[1:], "0", "")

		se += "00000"

		return s[:1] + se[1:4]
	})
}
