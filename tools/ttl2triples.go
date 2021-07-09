package main

import (
	"github.com/deiu/rdf2go"
	"github.com/pebbe/util"

	"fmt"
	"os"
)

var (
	x = util.CheckErr
)

func main() {

	fp, err := os.Open(os.Args[1])
	x(err)

	baseUri := "https://noordergraf.rug.nl/"

	// Create a new graph
	g := rdf2go.NewGraph(baseUri)
	x(err)

	// r is of type io.Reader
	g.Parse(fp, "text/turtle")

	fp.Close()

	for tri := range g.IterTriples() {
		fmt.Println(tri)
	}

}
