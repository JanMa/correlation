# Correlation [![GoDoc](https://godoc.org/gitlab.com/JanMa/correlation?status.svg)](https://godoc.org/gitlab.com/JanMa/correlation)

Correlation is a HTTP middleware for Go that adds correlation ids to incoming requests. It can be used with a standard net/http [Handler](https://golang.org/pkg/net/http/#Handler) and also be integrated with [Negroni](https://github.com/urfave/negroni).

## Usage

```go
 package main

  import (
      "net/http"

      "gitlab.com/JanMa/correlation"
  )

  var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      w.Write([]byte("hello world"))
  })

  func main() {
      correlationMiddleware := correlation.New(correlation.Options{
          CorrelationIDType: correlation.UUID,
      })

      http.ListenAndServe(":8080", correlationMiddleware.Handler(myHandler))
  }
```

The above example will add a `X-Correlation-ID` header to each incoming request containing a random UUID. 

### Available options

```go
c := correlation.New(correlation.Options{
    // HeaderName the name of the header to be used as correlation id. Defaults to `X-Correlation-ID`.
	CorrelationHeaderName: "X-Correlation-ID",
	// IDType the type of correlation id to generate. Defaults to `correlation.UUID`.
	CorrelationIDType: correlation.UUID,
	// CustomString the value to use when using a custom correlation id with the type correlation.Custom. Default is empty.
	CorrelationCustomString: "",
})
```

## Integration example

```go
package main

import (
    "net/http"

    "github.com/urfave/negroni"
    "gitlab.com/JanMa/correlation"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        w.Write([]byte("hello world!"))
    })

    correlationMiddleware := correlation.New(correlation.Options{
        CorrelationHeaderName: "Correlation-ID",
    })

    n := negroni.Classic()
    n.Use(negroni.HandlerFunc(correlationMiddleware.HandlerFuncWithNext))
    n.UseHandler(mux)

    n.Run(":8080")
}
```
