// Command gosp compiles .gosp template files to Go HTTP handlers.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jihoo12/gosp"
)

func main() {
	var source, output, pkg string
	flag.StringVar(&source, "src", ".", "directory containing .gosp files")
	flag.StringVar(&output, "out", "pages_gen.go", "generated .go output file")
	flag.StringVar(&pkg, "pkg", "main", "package for the generated handlers")
	flag.Parse()
	generated, err := gosp.Generate(source, pkg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gosp:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(output, generated, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "gosp:", err)
		os.Exit(1)
	}
}
