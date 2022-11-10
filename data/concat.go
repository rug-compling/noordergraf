package main

import (
	"github.com/pebbe/util"
	"github.com/rug-compling/noordergraf/go/nlsoundex"

	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	opt_n    = flag.Bool("n", false, "omit graph names")
	re       = regexp.MustCompile(`:fullName[ \t]+".*"`)
	reSUBDIR = regexp.MustCompile(`^[0-9][0-9][0-9]$`)
	reSUBPAT = regexp.MustCompile(`/[0-9][0-9][0-9]/`)
	x        = util.CheckErr
)

func main() {

	flag.Parse()

	x(os.Chdir("/net/noordergraf/data"))

	fp, err := os.Open("prefix.ttl")
	x(err)
	_, err = io.Copy(os.Stdout, fp)
	x(err)
	fp.Close()

	fmt.Println()

	doFile("ns.ttl")
	for _, dir := range []string{"bible", "item", "place", "site", "symbol", "todo"} {
		doDir(dir)
	}

	if !*opt_n {
		fmt.Println("<meta> {")
	}
	fmt.Printf("<https://noordergraf.rug.nl/> :dcModified \"%s+00:00\"^^xsd:dateTime .\n", time.Now().UTC().Format("2006-01-02T15:04:05"))
	if !*opt_n {
		fmt.Println("}")
	}

}

func doDir(dir string) {
	files, err := ioutil.ReadDir(dir)
	x(err)
	for _, file := range files {
		name := file.Name()
		if reSUBDIR.MatchString(name) {
			doDir(dir + "/" + name)
		} else {
			doFile(dir + "/" + name)
		}
	}
}

func doFile(filename string) {
	if !strings.HasSuffix(filename, ".ttl") {
		return
	}

	prefix := ":"
	if filename != "ns.ttl" {
		f := fix(filename)
		prefix = strings.Replace(f[:len(f)-4], "/", ":", -1)
	}

	inGraph := false
	doTime := true
	fp, err := os.Open(filename)
	x(err)
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		l := strings.TrimSpace(line)
		if l == "" {
			fmt.Println()
		} else if strings.HasPrefix(l, "@prefix") {
			if inGraph {
				inGraph = false
				fmt.Println("}")
			}
			fmt.Println(line)
		} else {
			if !inGraph {
				if !*opt_n {
					f := fix(filename)
					fmt.Printf("<%s> {\n", fix(f[:len(f)-4]))
					inGraph = true
				}
			}
			if doTime {
				fi, err := os.Stat(filename)
				x(err)
				fmt.Printf("%s :dcModified \"%s+00:00\"^^xsd:dateTime .\n", prefix, fi.ModTime().UTC().Format("2006-01-02T15:04:05"))
				doTime = false
			}
			fmt.Println(re.ReplaceAllStringFunc(line,
				func(s string) string {
					p1 := strings.Index(s, `"`) + 1
					p2 := strings.LastIndex(s, `"`)
					return fmt.Sprintf(`%s ; :fullNameSE "%s"`, s, nlsoundex.Soundex(s[p1:p2]))
				}))
		}
	}
	x(scanner.Err())

	if inGraph {
		fmt.Println("}")
	}
	fmt.Println()

}

func fix(s string) string {
	return reSUBPAT.ReplaceAllLiteralString(s, "/")
}
