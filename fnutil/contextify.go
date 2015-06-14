package fnutil

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
)

type key int

var fnutilUrlParamsKey key = 0

// ContextifyNetHttp creates a net/http wrapper around
// contextified handler.
func ContextifyNetHttp(
	c func(context.Context, http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		c(ctx, w, r)
	}
}

// ContextifyHttprouter creates a julienschmidt/httprouter wrapper
// around contextified handler. Url params from httprouter are
// included into the context passed to the handler.
func ContextifyHttprouter(
	c func(context.Context, http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(context.Background(), fnutilUrlParamsKey, ps)
		c(ctx, w, r)
	}
}

// GetUrlParams extracts url params from context if there are any.
func GetUrlParams(c context.Context) (params httprouter.Params, ok bool) {
	params, ok = c.Value(fnutilUrlParamsKey).(httprouter.Params)
	return params, ok
}

// ContextifyHandleFunc wraps a net/http compatible HandleFunc
// to match the context aware signature.
func ContextifyHandleFunc(
	netHttpHandleFunc func(http.ResponseWriter, *http.Request),
) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		netHttpHandleFunc(w, r)
	}
}

// ContextifyHandler wraps a net/http Handler interface to match
// context aware signature.
func ContextifyHandler(
	netHttpHandler http.Handler,
) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		netHttpHandler.ServeHTTP(w, r)
	}
}

// ContextifyMiddleware wraps a net/http compatible middleware
// decorator to match the context aware signature. Context bypasses
// the net/http middleware and is passed to next layer without
// modifications.
func ContextifyMiddleware(
	netHttpMiddleware func(func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request),
) func(
	func(context.Context, http.ResponseWriter, *http.Request),
) func(context.Context, http.ResponseWriter, *http.Request) {

	return func(
		f func(context.Context, http.ResponseWriter, *http.Request),
	) func(context.Context, http.ResponseWriter, *http.Request) {

		return func(c context.Context, w http.ResponseWriter, r *http.Request) {
			var wrapped = func(_w http.ResponseWriter, _r *http.Request) {
				f(c, _w, _r)
			}
			netHttpMiddleware(wrapped)
		}
	}
}
