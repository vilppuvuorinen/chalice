# Chalice

Chalice provides a convinient way to chain your HTTP middleware
*functions* inspired by justinas/alice.

Alice works around http.Handler interface implementations whereas
Chalice works strictly around functions matching x/net/context aware
ServeHTTP signature:

```
func(context.Context, http.ResponseWriter, *http.Request)
```

There are also several wrapper functions that allow usage of:

* http.Handler implementations
* plain ServeHTTP implementations
* possibly other handler in the future

Adapters are also provided for the other end of the chain allowing usage
of the resulting handler chain with:

* net/http
* julienschmidt/httprouter (hard dependecy in url params)
* possibly others in the future

## TODO

* Short usage example in README
* Tests
* Prettifying the code
