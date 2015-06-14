package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"github.com/vilppuvuorinen/chalice/fnutil"
	_m "github.com/vilppuvuorinen/chalice/middleware"
)

func IndexHandle(c context.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func SecretHandle(c context.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Secret")
}

func PanicHandle(c context.Context, w http.ResponseWriter, r *http.Request) {
	panic("Panic triggered")
}

func main() {
	log.Println("Starting server...")

	var baseHandle = fnutil.MkPartial(
		_m.Logger,
		_m.PanicRecover,
	)

	r := httprouter.New()

	r.NotFound = fnutil.ContextifyNetHttp(
		baseHandle(
			fnutil.ContextifyHandler(http.NotFound),
		),
	)

	// Index
	r.GET("/", fnutil.ContextifyHttprouter(
		baseHandle(
			IndexHandle,
		)),
	)
	// Login with always failing authentication
	r.GET("/login", fnutil.ContextifyHttprouter(
		baseHandle(
			SecretHandle,
			_m.NewBasicAuth(func(username, password string) bool {
				return false
			}),
		)),
	)

	// Panicer
	r.GET("/panic", fnutil.ContextifyHttprouter(
		baseHandle(
			PanicHandle,
		)),
	)

	log.Fatal(http.ListenAndServe(":8000", r))
}
