package main

import (
	"github.com/pebbe/util"
	"github.com/rug-compling/noordergraf/go/nlsoundex"

	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	x  = util.CheckErr
	re = regexp.MustCompile(`:fullname[ \t]+".*"[ \t]*;[ \t]*$`)
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(re.ReplaceAllStringFunc(scanner.Text(),
			func(s string) string {
				p1 := strings.Index(s, `"`) + 1
				p2 := strings.LastIndex(s, `"`)
				return fmt.Sprintf(`%s :fullnameSE "%s" ;`, s, nlsoundex.Soundex(s[p1:p2]))
			}))
	}
	x(scanner.Err())
}
