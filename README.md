go-deps
=======

[![xrefs](https://sourcegraph.com/api/repos/github.com/sourcegraph/go-deps/badges/xrefs.png)](https://sourcegraph.com/github.com/sourcegraph/go-deps)
[![funcs](https://sourcegraph.com/api/repos/github.com/sourcegraph/go-deps/badges/funcs.png)](https://sourcegraph.com/github.com/sourcegraph/go-deps)
[![top func](https://sourcegraph.com/api/repos/github.com/sourcegraph/go-deps/badges/top-func.png)](https://sourcegraph.com/github.com/sourcegraph/go-deps)
[![library users](https://sourcegraph.com/api/repos/github.com/sourcegraph/go-deps/badges/library-users.png)](https://sourcegraph.com/github.com/sourcegraph/go-deps)
[![status](https://sourcegraph.com/api/repos/github.com/sourcegraph/go-deps/badges/status.png)](https://sourcegraph.com/github.com/sourcegraph/go-deps)

Package deps analyzes and recursively installs Go package dependencies. It is library functionality similar to `go get`.

Docs: [go-deps on Sourcegraph](https://sourcegraph.com/github.com/sourcegraph/go-deps)

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
