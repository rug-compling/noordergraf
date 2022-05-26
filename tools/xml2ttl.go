package main

import (
	"github.com/pebbe/util"

	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	x     = util.CheckErr
	reIRI = regexp.MustCompile(`<.*?>`)
)

func main() {
	b, err := ioutil.ReadFile(os.Args[1])
	x(err)

	cmd := exec.Command("rapper", strings.Fields("-i rdfxml -o turtle -I https://noordergraf.rug.nl/ -q -")...)
	stdin, err := cmd.StdinPipe()
	x(err)
	go func() {
		defer stdin.Close()
		stdin.Write(b)
	}()

	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	x(err)

	prefix := make(map[string]string)
	var base string
	lines := strings.SplitAfter(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "@base") {
			base = strings.Fields(line)[1]
			base = base[1 : len(base)-1]
			continue
		}
		if strings.HasPrefix(line, "@prefix") {
			aa := strings.Fields(line)
			if strings.HasPrefix(aa[2], "<item/") && aa[1] != "item:" {
				line = strings.Replace(line, "<", "<"+base, 1)
				fmt.Print(line)
			} else {
				prefix[aa[2][1:len(aa[2])-1]] = aa[1]
			}
			continue
		}
		line = strings.Replace(line, "    ", "  ", -1)
		line = reIRI.ReplaceAllStringFunc(line, func(s string) string {
			iri := s[1 : len(s)-1]
			for key, value := range prefix {
				if strings.HasPrefix(iri, key) {
					return value + iri[len(key):]
				}
			}
			return s
		})
		fmt.Print(line)
	}
}
