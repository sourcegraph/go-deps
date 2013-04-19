go-deps
=======
Package deps analyzes and recursively installs Go packages. It is library functionality similar to `go get`.

Source: [github.com/sqs/go-deps](https://github.com/sqs/go-deps)
Docs: [godoc.org/github.com/sqs/go-deps](http://godoc.org/github.com/sqs/go-deps)

## Installation

	go get github.com/sqs/go-deps

## Example Usage

    import (
        "github.com/sqs/go-deps"
    )

	pkg, _ := deps.Read(test.importPath)
    for _, p := range pkg.DepsNotFound {
        p.Download()
    }

## Authors

* Quinn Slack <qslack@qslack.com>
