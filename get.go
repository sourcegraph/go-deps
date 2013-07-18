package deps

import (
	"fmt"
	"os"
	"os/exec"
)

// Downloads (but does not install) the package at the given import path to the first GOPATH tree.
// Works like `go get -d`.
func (c *Context) Download(importPath string) error {
	return c.GoGet(importPath, DownloadOnly)
}

type GetMode int

const (
	DownloadOnly GetMode = iota
	Update
	Verbose
)

// Runs `go get` on importPath with the options specified in mode.
func (c *Context) GoGet(importPath string, mode GetMode) error {
	args := []string{"get"}
	if mode&DownloadOnly > 0 {
		args = append(args, "-d")
	}
	if mode&Update > 0 {
		args = append(args, "-u")
	}
	if mode&Verbose > 0 {
		args = append(args, "-x", "-v")
	}
	args = append(args, "--", string(importPath))
	return c.gocmd(args)
}

func (c *Context) gocmd(args []string) error {
	fmt.Fprintf(os.Stderr, "go %v\n", args)
	cmd := exec.Command("go", args...)
	cmd.Env = []string{"GOROOT=" + c.GOROOT, "GOPATH=" + c.GOPATH, "GIT_ASKPASS=echo", "PATH=" + os.Getenv("PATH")}
	cmd.Stdout = c.Out
	cmd.Stderr = c.Err
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error: `go get` with args %v: %s", args, err)
	}
	return nil
}

// Calls Download with the go/build.Default build context.
func Download(importPath string) error {
	return Default.Download(importPath)
}

// Calls GoGet with the go/build.Default build context.
func GoGet(importPath string, mode GetMode) error {
	return Default.GoGet(importPath, mode)
}
