# Kocha-urlrouter [![Build Status](https://travis-ci.org/naoina/kocha-urlrouter.png?branch=master)](https://travis-ci.org/naoina/kocha-urlrouter)

Better URL router collection for [Go](http://golang.org)

## Installation

Interface:

    go get -u github.com/naoina/kocha-urlrouter

Implementation:

    go get -u github.com/naoina/kocha-urlrouter/doublearray

Kocha-urlrouter has multiple URL router implementations. See [Implementations](#implementations).

## Usage

```go
package main

import (
    "github.com/naoina/kocha-urlrouter"
    _ "github.com/naoina/kocha-urlrouter/doublearray"
)

type route struct {
    name string
}

func main() {
    router, err := urlrouter.NewURLRouter("doublearray")
    if err != nil {
        panic(err)
    }
    router.Build([]*urlrouter.Record{
        urlrouter.NewRecord("/", &route{"root"}),
        urlrouter.NewRecord("/user/:id", &route{"user"}),
        urlrouter.NewRecord("/static/*filepath", &route{"static"}),
    })

    router.Lookup("/")                    // returns *route{"root"}, nil map
    router.Lookup("/user/hoge")           // returns *route{"user"}, map of {"id": "hoge"}
    router.Lookup("/static/path/to/file") // returns *route{"static"}, map of {"filepath": "path/to/file"}
}
```

See [Godoc](http://godoc.org/github.com/naoina/kocha-urlrouter) for more docs.

## Implementations

* Double-Array `github.com/naoina/kocha-urlrouter/doublearray`
* Regular-Expression `github.com/naoina/kocha-urlrouter/regexp`
* Ternary Search Tree `github.com/naoina/kocha-urlrouter/tst`

## Benchmark

    cd $GOPATH/github.com/naoina/kocha-urlrouter
    go test -bench . -benchmem ./...

## License

Kocha-urlrouter is licensed under the MIT
