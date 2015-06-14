package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vilppuvuorinen/chalice"
	"github.com/vilppuvuorinen/chalice/compat/httproutercmpt"
	_m "github.com/vilppuvuorinen/chalice/middleware"
	"golang.org/x/net/context"
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

	var baseHandle = chalice.MkPartial(
		_m.Logger,
		_m.PanicRecover,
	)

	r := httprouter.New()

	r.NotFound = chalice.CallWithContext(baseHandle(
		chalice.ContextifyHandleFunc(http.NotFound),
	))

	// Index
	r.GET(
		"/",
		httproutercmpt.CallWithContext(baseHandle(IndexHandle)),
	)

	// Login with always failing authentication
	r.GET(
		"/login",
		httproutercmpt.CallWithContext(baseHandle(
			SecretHandle,
			_m.NewBasicAuth(func(username, password string) bool { return false }),
		)),
	)

	// Panicer
	r.GET(
		"/panic",
		httproutercmpt.CallWithContext(baseHandle(PanicHandle)),
	)

	log.Fatal(http.ListenAndServe(":8000", r))
}
