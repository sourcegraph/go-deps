go-deps
=======
Package deps analyzes and recursively installs Go package dependencies. It is library functionality similar to `go get`.

Docs: [go-deps on Sourcegraph](https://sourcegraph.com/repos/github.com/sourcegraph/go-deps)

## Installation

	go get github.com/sourcegraph/go-deps

## Example Usage

    import (
        "github.com/sourcegraph/go-deps"
    )

	pkg, _ := deps.Read(test.importPath)
    for _, p := range pkg.DepsNotFound {
        p.Download()
    }

## Authors

* Quinn Slack <qslack@qslack.com>
