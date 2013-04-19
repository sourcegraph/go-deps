package deps

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
)

// Downloads the package at the given import path to the first GOPATH
// tree. Works like `go get -d`.
func (p ImportPath) Download() error {
	cmd := exec.Command("go", "get", "-d", "--", string(p))
	cmd.Env = []string{"GOPATH=" + build.Default.GOPATH, "GIT_ASKPASS=echo", "PATH=" + os.Getenv("PATH")}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %s", err, string(out))
	}
	return nil
}
