package main

/*

Dit programma controleert of bestanden onder data/item en www/img in
de juiste subdirectory staan.

*/

import (
	"github.com/pebbe/util"

	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	re = regexp.MustCompile(`^[0-9][0-9][0-9]$`)
	x  = util.CheckErr
)

func main() {
	chkDir("/net/noordergraf/data/item")
	chkDir("/net/noordergraf/www/img")
}

func chkDir(dirname string) {
	dirs, err := ioutil.ReadDir(dirname)
	x(err)
	for _, dir := range dirs {
		if name := dir.Name(); re.MatchString(name) {
			chkFiles(dirname, name)
		}
	}
}

func chkFiles(dir1, dir2 string) {
	files, err := ioutil.ReadDir(dir1 + "/" + dir2)
	x(err)
	for _, file := range files {
		name := file.Name()
		if i := strings.LastIndex(name, "."); i > 0 {
			name = name[:i]
		}
		if !strings.HasSuffix(name, dir2) {
			x(fmt.Errorf("File %s in wrong directory %s", file.Name(), dir1+"/"+dir2))
		}
	}
}
