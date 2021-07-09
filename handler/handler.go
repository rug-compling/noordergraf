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
	fmt.Println(err, string(b))

	env := os.Environ()
	sort.Strings(env)
	fmt.Println(strings.Join(env, "\n"))
}
