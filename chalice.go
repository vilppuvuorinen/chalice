// Package chalice provides a set of functions for developing
// functional-ish REST API. Also there is the ulterior motive
// of pushing godoc to the limits with borderline incomprehensible
// func signatures.
package chalice

import (
	"golang.org/x/net/context"
	"net/http"
)

// MkHandleFunc creates a context aware http handler function
// with defined middleware and handler function.
func MkHandleFunc(
	handler func(context.Context, http.ResponseWriter, *http.Request),
	middleware ...func(func(context.Context, http.ResponseWriter, *http.Request)) func(context.Context, http.ResponseWriter, *http.Request),
) func(context.Context, http.ResponseWriter, *http.Request) {

	end := len(middleware) - 1
	h := handler
	for i := end; i >= 0; i-- {
		h = middleware[i](h)
	}

	return h
}

// MkPartial creates a partial MkHandleFunc that accepts more
// middleware and a handler.
func MkPartial(
	base ...(func(func(context.Context, http.ResponseWriter, *http.Request)) func(context.Context, http.ResponseWriter, *http.Request)),
) func(
	func(context.Context, http.ResponseWriter, *http.Request),
	...(func(func(context.Context, http.ResponseWriter, *http.Request),
	) func(context.Context, http.ResponseWriter, *http.Request)),
) func(context.Context, http.ResponseWriter, *http.Request) {

	return func(
		handler func(context.Context, http.ResponseWriter, *http.Request),
		middleware ...func(func(context.Context, http.ResponseWriter, *http.Request)) func(context.Context, http.ResponseWriter, *http.Request),
	) func(context.Context, http.ResponseWriter, *http.Request) {
		return MkHandleFunc(handler, append(base, middleware...)...)
	}
}
