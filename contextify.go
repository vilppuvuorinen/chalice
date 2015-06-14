package chalice

import (
	"golang.org/x/net/context"
	"net/http"
)

type key int

var URLParamsKey key = 1

// Param is a single URL parameter adopted from httprouter.
type Param struct {
	Key, Value string
}

// Params is a Param-slice constructed from data from router or
// middleware.
type Params []Param

func (p Params) ByName(name string) (value string, ok bool) {
	value = ""
	ok = false
	for i := range p {
		if p[i].Key == name {
			value = p[i].Value
			ok = true
			break
		}
	}
	return value, ok
}

// CallWithContext creates a net/http wrapper around
// contextified handler.
func CallWithContext(c func(context.Context, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		c(ctx, w, r)
	}
}

// GetURLParams extracts url params from context if there are any.
func GetURLParams(c context.Context) (params Params, ok bool) {
	params, ok = c.Value(URLParamsKey).(Params)
	return params, ok
}

// SetURLParams stores Params to context.
func SetURLParams(c context.Context, p Params) context.Context {
	return context.WithValue(c, URLParamsKey, p)
}

// ContextifyHandleFunc wraps a net/http compatible HandleFunc
// to match the context aware signature.
func ContextifyHandleFunc(netHTTPHandleFunc func(http.ResponseWriter, *http.Request)) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		netHTTPHandleFunc(w, r)
	}
}

// ContextifyHandler wraps a net/http Handler interface to match
// context aware signature.
func ContextifyHandler(netHTTPHandler http.Handler) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		netHTTPHandler.ServeHTTP(w, r)
	}
}

// ContextifyMiddleware wraps a net/http compatible middleware
// decorator to match the context aware signature. Context bypasses
// the net/http middleware and is passed to next layer without
// modifications.
func ContextifyMiddleware(
	netHTTPMiddleware func(func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request),
) func(func(context.Context, http.ResponseWriter, *http.Request)) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(f func(context.Context, http.ResponseWriter, *http.Request)) func(context.Context, http.ResponseWriter, *http.Request) {
		return func(c context.Context, w http.ResponseWriter, r *http.Request) {
			var wrapped = func(_w http.ResponseWriter, _r *http.Request) {
				f(c, _w, _r)
			}
			netHTTPMiddleware(wrapped)
		}
	}
}
