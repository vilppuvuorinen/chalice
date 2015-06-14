# Chalice

[![GoDoc](https://godoc.org/github.com/vilppuvuorinen/chalice?status.svg)](https://godoc.org/github.com/vilppuvuorinen/chalice)

Chalice provides a convinient way to chain your HTTP middleware
*functions* inspired by justinas/alice.

Alice works around http.Handler interface implementations whereas
Chalice works strictly around functions matching x/net/context aware
ServeHTTP signature:

```go
func(context.Context, http.ResponseWriter, *http.Request)
```

There are also several wrapper functions that allow usage of:

* http.Handler implementations
* plain ServeHTTP implementations
* possibly other handler in the future

Adapters are also provided for the other end of the chain allowing usage
of the resulting handler chain with:

* net/http
* julienschmidt/httprouter
* possibly others in the future

## Why?

Alice is a great tool but basic net/http lacks global-less context.
Passing variables through handler stack requires breaking basic net/http
signature. In addition, the borderline incomprehensible function
signatures were also a driving factor.

## Usage

Arbitrary routers can be used by wrapping context aware handle
functions in compatibility decorator.

```go
func IndexHandle(c context.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

r := httprouter.New()
r.GET(
  "/",
  httproutercmpt.CallWithContext(IndexHandle),
)

http.ListenAndServe(":8000", r)
```

Middleware can be chained using  _MkHandleFunc_ to chain middleware
just like Alice's _New_ function.

```go
r.GET(
  "/",
  httproutercmpt.CallWithContext(
    MkHandleFunc(
      IndexHandle,
      middleware.Logger,
      middleware.PanicRecover,
    ),
  ),
)
```

Default chain can be constructed as a partial _MkHandleFunc_ using
_MkPartial_.

```go
var baseHandle = chalice.MkPartial(
  middleware.Logger,
  middleware.PanicRecover,
)

r.GET(
  "/",
  httproutercmpt.CallWithContext(
    baseHandle(
      IndexHandle,
      middleware.Logger,
      middleware.PanicRecover,
      customMiddleware,
    ),
  ),
)
```

Compatibility wrappers allow using basic net/http HandleFuncs as context
aware HandleFuncs.

```go
r := httprouter.New()
r.NotFound = chalice.CallWithContext(baseHandle(
  chalice.ContextifyHandleFunc(http.NotFound),
))
```

See full [example](example/example.go).

## TODO

* Tests
