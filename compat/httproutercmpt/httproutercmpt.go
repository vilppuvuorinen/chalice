// Package httprouter provides Contextify wrapper for
// httprouter specific HandleFunc implementations.
package httproutercmpt

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vilppuvuorinen/chalice"
	"golang.org/x/net/context"
)

// CallWithContext creates a julienschmidt/httprouter wrapper
// around contextified handler. Url params from httprouter are
// included into the context passed to the handler.
func CallWithContext(
	c func(context.Context, http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(
			context.Background(),
			chalice.URLParamsKey,
			translateParams(ps),
		)
		c(ctx, w, r)
	}
}

func translateParams(urlParams httprouter.Params) chalice.Params {
	if urlParams == nil {
		return nil
	}

	params := make(chalice.Params, len(urlParams))

	for i, p := range urlParams {
		params[i] = chalice.Param{p.Key, p.Value}
	}
	return params
}
