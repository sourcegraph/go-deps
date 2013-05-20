package deps

import (
	"fmt"
	"os"
	"os/exec"
)

// Downloads the package at the given import path to the first GOPATH
// tree. Works like `go get -d`.
func (c *Context) Download(importPath ImportPath) error {
	cmd := exec.Command("go", "get", "-d", "--", string(importPath))
	cmd.Env = []string{"GOPATH=" + c.GOPATH, "GIT_ASKPASS=echo", "PATH=" + os.Getenv("PATH")}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("download %q: %s: %s", importPath, err, string(out))
	}
	return nil
}

// Calls Download with the go/build.Default build context.
func Download(importPath ImportPath) error {
	return Default.Download(importPath)
}
