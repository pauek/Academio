
package main

import (
	"Academio/content"
	"fmt"
	"flag"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: dir2id <dir>")
	} else {
		fmt.Printf("%s\n", content.ToID(args[0]))
	}
}