# midway

[![GoDoc](https://godoc.org/github.com/novrin/midway?status.svg)](https://pkg.go.dev/github.com/novrin/midway) 
![tests](https://github.com/novrin/midway/workflows/tests/badge.svg) 
[![Go Report Card](https://goreportcard.com/badge/github.com/novrin/midway)](https://goreportcard.com/report/github.com/novrin/midway)

`midway` is a micro Go package for arranging your HTTP middleware.

### Installation

```shell
go get github.com/novrin/midway
``` 

## Usage

A `Queue` returns a `Middleware` where the slice of given middlewares are applied first-in-first-out. The last middleware in the slice will execute first.

```go
queued := midway.Queue(corsHeaders, secureHeaders)
http.ListenAndServe(":1313", queued(app))
// secureHeaders(corsHeaders(app))
```

A `Stack` returns a `Middleware` where the slice of given middlewares are applied last-in-first-out. The first middleware in the slice will execute first.

```go
stacked := midway.Stack(corsHeaders, secureHeaders)
http.ListenAndServe(":1313", queued(app))
// corsHeaders(secureHeaders(app))
```

A complete example using a `Queue`.

```go
package main

import (
	"net/http"

	"github.com/novrin/midway"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func secureHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		...
		h.ServeHTTP(w, r)
	})
}

func corsHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		...
		h.ServeHTTP(w, r)
	})
}

func main() {
	// Use Queue to arrange middleware to execute secureHeaders first.
	queued := midway.Queue(corsHeaders, secureHeaders)

	app := http.HandlerFunc(hello)
	http.ListenAndServe(":1313", queued(app))
}
```

## License

Copyright (c) 2023-present [novrin](https://github.com/novrin)

Licensed under [MIT License](./LICENSE)