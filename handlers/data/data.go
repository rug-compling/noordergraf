package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func main() {
	fmt.Print("Content-type: text/plain\n\n")
	b, err := exec.Command("id").CombinedOutput()
	if err == nil {
		fmt.Println(string(b))
	} else {
		fmt.Println(err)
		fmt.Println()
	}

	env := os.Environ()
	sort.Strings(env)
	fmt.Println(strings.Join(env, "\n"))
}
