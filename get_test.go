package deps

import (
	"go/build"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	ctx := Context(build.Default)
	ctx.GOPATH, _ = filepath.Abs("testdata")
	tests := []struct {
		importPath string
		pkg        Package
	}{
		{"imports_stdlib", Package{
			ImportPath: "imports_stdlib",
			Name:       "imports_stdlib",
			Goroot:     false, Standard: false,
			Imports:    []string{"errors"},
			Deps:       []string{"errors", "runtime", "unsafe"},
			Incomplete: false,
			Error:      nil,
			DepsErrors: nil,
		}},
		{"imports_github", Package{
			ImportPath: "imports_github",
			Name:       "imports_github",
			Goroot:     false, Standard: false,
			Imports:      []string{"github.com/sqs/go-deps"},
			Deps:         []string{"github.com/sqs/go-deps", "runtime", "unsafe"},
			DepsNotFound: []string{"github.com/sqs/go-deps"},
			Incomplete:   true,
			Error:        nil,
			DepsErrors: []*PackageError{
				{
					ImportStack: []string{"imports_github", "github.com/sqs/go-deps"},
					Pos:         "testdata/src/imports_github/imports_github.go:3:8",
					Err:         "cannot find package \"github.com/sqs/go-deps\"",
				},
			},
		}},
		{"imports_source_pkg", Package{
			ImportPath: "imports_source_pkg",
			Name:       "imports_source_pkg",
			Goroot:     false, Standard: false,
			Imports:    []string{"src_pkg"},
			Deps:       []string{"runtime", "src_pkg", "unsafe"},
			Incomplete: false,
			Error:      nil,
			DepsErrors: nil,
		}},
		{"missing_import", Package{
			ImportPath: "missing_import",
			Name:       "missing_import",
			Goroot:     false, Standard: false,
			Imports:      []string{"doesnotexist", "github.com/example/alsodoesntexist"},
			Deps:         []string{"doesnotexist", "github.com/example/alsodoesntexist", "runtime", "unsafe"},
			DepsNotFound: []string{"doesnotexist", "github.com/example/alsodoesntexist"},
			Incomplete:   true,
			Error:        nil,
			DepsErrors: []*PackageError{
				{
					ImportStack: []string{"missing_import", "doesnotexist"},
					Pos:         "testdata/src/missing_import/missing_import.go:4:2",
					Err:         "cannot find package \"doesnotexist\"",
				},
				{
					ImportStack: []string{"missing_import", "github.com/example/alsodoesntexist"},
					Pos:         "testdata/src/missing_import/missing_import.go:5:2",
					Err:         "cannot find package \"github.com/example/alsodoesntexist\"",
				},
			},
		}},
		{"github.com/example/doesntexist", Package{
			ImportPath: "github.com/example/doesntexist",
			Name:       "",
			Goroot:     false, Standard: false,
			Imports:    nil,
			Deps:       nil,
			Incomplete: true,
			Error: &PackageError{
				ImportStack: []string{"github.com/example/doesntexist"},
				Pos:         "",
				Err:         "cannot find package \"github.com/example/doesntexist\"",
			},
			DepsErrors: nil,
		}},
	}

	for _, test := range tests {
		pkg, err := ctx.Read(test.importPath)
		if err != nil {
			t.Fatalf("%s: Read failed: %s", test.importPath, err)
		}
		checkPackagesEqual(t, test.importPath, &test.pkg, pkg)
	}
}

// Tests two Packages for equality, checking only those fields that we care
// about in tests. E.g., ignores fields that contain absolute paths
// (for now).
func checkPackagesEqual(t *testing.T, importPath string, exp, actual *Package) {
	if exp.ImportPath != actual.ImportPath {
		t.Errorf("%s: ImportPath: want %v, got %v", importPath, exp.ImportPath, actual.ImportPath)
	}
	if exp.Name != actual.Name {
		t.Errorf("%s: Name: want %v, got %v", importPath, exp.Name, actual.Name)
	}
	if exp.Goroot != actual.Goroot {
		t.Errorf("%s: Goroot: want %v, got %v", importPath, exp.Goroot, actual.Goroot)
	}
	if exp.Standard != actual.Standard {
		t.Errorf("%s: Standard: want %v, got %v", importPath, exp.Standard, actual.Standard)
	}
	if !reflect.DeepEqual(exp.Imports, actual.Imports) {
		t.Errorf("%s: Imports: want %v, got %v", importPath, exp.Imports, actual.Imports)
	}
	if !reflect.DeepEqual(exp.Deps, actual.Deps) {
		t.Errorf("%s: Deps: want %v, got %v", importPath, exp.Deps, actual.Deps)
	}
	if !reflect.DeepEqual(exp.DepsNotFound, actual.DepsNotFound) {
		t.Errorf("%s: DepsNotFound: want %v, got %v", importPath, exp.DepsNotFound, actual.DepsNotFound)
	}
	if exp.Incomplete != actual.Incomplete {
		t.Errorf("%s: Incomplete: want %v, got %v", importPath, exp.Incomplete, actual.Incomplete)
	}
	checkPackageErrorsEqual(t, importPath, exp.Error, actual.Error)
	if len(exp.DepsErrors) == len(actual.DepsErrors) {
		for i, expe := range exp.DepsErrors {
			actuale := actual.DepsErrors[i]
			checkPackageErrorsEqual(t, importPath, expe, actuale)
		}
	} else {
		t.Errorf("%s: DepsErrors: want len %d (%v), got len %d (%v)", importPath, len(exp.DepsErrors), exp.DepsErrors, len(actual.DepsErrors), actual.DepsErrors)
	}
}

func derefPackageErrors(errps []*PackageError) (errs []PackageError) {
	errs = make([]PackageError, len(errps))
	for i, errp := range errps {
		errs[i] = *errp
	}
	return
}

// Tests two PackageErrors for equality, ignoring absolute paths.
func checkPackageErrorsEqual(t *testing.T, importPath string, exp, actual *PackageError) {
	if exp == actual {
		return
	}
	if exp == nil || actual == nil {
		t.Errorf("%s DepsErrors: want %v, got %v", importPath, exp, actual)
		return
	}
	if !reflect.DeepEqual(exp.ImportStack, actual.ImportStack) {
		t.Errorf("%s DepsErrors.ImportStack: want %#v, got %#v", importPath, exp.ImportStack, actual.ImportStack)
	}
	if exp.Pos != actual.Pos {
		t.Errorf("%s DepsErrors.Pos: want %#v, got %#v", importPath, exp.Pos, actual.Pos)
	}
	if !strings.HasPrefix(actual.Err, exp.Err) {
		t.Errorf("%s DepsErrors.Err: want prefix %#v, got %#v", importPath, exp.Err, actual.Err)
	}
}
