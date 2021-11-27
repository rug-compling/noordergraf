package main

import (
	"github.com/pebbe/util"
	"github.com/rug-compling/noordergraf/go/nlsoundex"

	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var (
	lines = make([]string, 0)

	x = util.CheckErr
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, nlsoundex.Soundex(line)+"\t"+line)
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
