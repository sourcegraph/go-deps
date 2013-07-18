package deps

import (
	"encoding/json"
	"go/build"
	"io"
	"os"
	"os/exec"
	"strings"
)

type bar int

// Relevant fields from Package struct in `go list` (cmd/go/pkg.go).
type Package struct {
	// Note: These fields are part of the go command's public API.
	// See list.go.  It is okay to add fields, but not to change or
	// remove existing ones.  Keep in sync with list.go
	Foo        bar
	Dir        string `json:",omitempty"` // directory containing package sources
	ImportPath string `json:",omitempty"` // import path of package in dir
	Name       string `json:",omitempty"` // package name
	Target     string `json:",omitempty"` // install path
	Goroot     bool   `json:",omitempty"` // is this package found in the Go root?
	Standard   bool   `json:",omitempty"` // is this package part of the standard Go library?
	Stale      bool   `json:",omitempty"` // would 'go install' do anything for this package?
	Root       string `json:",omitempty"` // Go root or Go path dir containing this package

	// Dependency information
	Imports      []string `json:",omitempty"` // import paths used by this package
	Deps         []string `json:",omitempty"` // all (recursively) imported dependencies
	DepsNotFound []string `json:",omitempty"` // all (recursive) deps that were not found

	// Error information
	Incomplete bool            `json:",omitempty"` // was there an error loading this package or dependencies?
	Error      *PackageError   `json:",omitempty"` // error loading this package (not dependencies)
	DepsErrors []*PackageError `json:",omitempty"` // errors loading dependencies
}

// A PackageError describes an error loading information about a package.
type PackageError struct {
	ImportStack []string // shortest path from package named on command line to this one
	Pos         string   // position of error
	Err         string   // the error itself
}

func (p *PackageError) Error() string {
	if p.Pos != "" {
		// Omit import stack.  The full path to the file where the error
		// is the most important thing.
		return p.Pos + ": " + p.Err
	}
	if len(p.ImportStack) == 0 {
		return p.Err
	}
	return "package " + strings.Join(p.ImportStack, "\n\timports ") + ": " + p.Err
}

type Context struct {
	build.Context
	Err io.Writer
	Out io.Writer
}

var Default Context = Context{
	Context: build.Default,
	Err:     os.Stderr,
	Out:     os.Stdout,
}

// Reads package info for the package at importPath from `go list -json`. If importPath is "",
// reads the package in the current directory.
func (c *Context) Read(importPath string) (pkg *Package, err error) {
	cmd := exec.Command("go", "list", "-e", "-json", importPath)
	cmd.Env = []string{"GOROOT=" + c.GOROOT, "GOPATH=" + c.GOPATH}
	var out []byte
	if out, err = cmd.Output(); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(out, &pkg); err != nil {
		return nil, err
	}

	for _, deperr := range pkg.DepsErrors {
		if strings.HasPrefix(deperr.Err, "cannot find package") {
			pkg.DepsNotFound = append(pkg.DepsNotFound, deperr.ImportStack[len(deperr.ImportStack)-1])
		}
	}

	return pkg, err
}

// Calls Read(importPath) with the go/build.Default build context.
func Read(importPath string) (pkg *Package, err error) {
	return Default.Read(importPath)
}
