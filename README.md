# Kocha-urlrouter [![Build Status](https://travis-ci.org/naoina/kocha-urlrouter.png?branch=master)](https://travis-ci.org/naoina/kocha-urlrouter)

Better URL router collection for [Go](http://golang.org)

**Note**:
Kocha-urlrouter will be used as a sandbox for some implementations of a URL router. If you want a fast URL router, Please use [Denco](https://github.com/naoina/denco) instead.

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
    router := urlrouter.NewURLRouter("doublearray")
    router.Build([]urlrouter.Record{
        urlrouter.NewRecord("/", &route{"root"}),
        urlrouter.NewRecord("/user/:id", &route{"user"}),
        urlrouter.NewRecord("/user/:name/:id", &route{"username"}),
        urlrouter.NewRecord("/static/*filepath", &route{"static"}),
    })

    router.Lookup("/")                    // returns *route{"root"}, nil slice.
    router.Lookup("/user/hoge")           // returns *route{"user"}, []urlrouter.Param{{"id", "hoge"}}
    router.Lookup("/user/hoge/7")           // returns *route{"username"}, []urlrouter.Param{{"name", "hoge"}, {"id", "7"}}
    router.Lookup("/static/path/to/file") // returns *route{"static"}, []urlrouter.Param{{"filepath", "path/to/file"}}
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
