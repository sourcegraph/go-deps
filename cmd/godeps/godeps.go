package main

import (
	"flag"
	"fmt"
	"github.com/sqs/go-deps"
	"os"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "must specify a single, locally installed package\n")
		os.Exit(1)
	}

	importPath := flag.Arg(0)
	pkg, err := deps.Read(importPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading package: %s\n", err)
		os.Exit(1)
	}

	if pkg.Dir == "" {
		// User specified package that's not installed locally.
		fmt.Fprintf(os.Stderr, "can't find package: import path %s\n", importPath)
		os.Exit(1)
	}

	fmt.Printf("# %s\n", importPath)
	for _, p := range pkg.DepsNotFound {
		fmt.Printf("%s\n", p)
		if err := deps.Download(p); err != nil {
			fmt.Fprintf(os.Stderr, "can't download dep: import path %s: %s\n", p, err)
			os.Exit(1)
		}
	}
}
