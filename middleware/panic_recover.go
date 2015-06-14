package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/vilppuvuorinen/chalice/httputil"
	"golang.org/x/net/context"
)

// PanicRecover recovers from panic and returns HTTP 500 if possible.
func PanicRecover(f func(
	context.Context,
	http.ResponseWriter,
	*http.Request,
)) func(context.Context, http.ResponseWriter, *http.Request) {

	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				httputil.Error(w, http.StatusInternalServerError)
			}
		}()

		f(c, w, r)
	}
}
