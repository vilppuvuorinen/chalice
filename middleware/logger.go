package middleware

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

// Logger logs incoming requests and their status codes.
// http.ResponseWriter is wrapped into a loggable interface
// implementation that keeps track of the resulting status
// code.
func Logger(
	f func(context.Context, http.ResponseWriter, *http.Request,
	)) func(context.Context, http.ResponseWriter, *http.Request) {

	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := newLoggedResponseWriter(w)

		f(c, lw, r)

		if lw.Status < 0 {
			lw.Status = http.StatusOK
		}

		log.Printf(
			"INFO:  %s %s %d %s %s",
			r.Method,
			r.RequestURI,
			lw.Status,
			http.StatusText(lw.Status),
			time.Since(start),
		)
	}
}

type loggedResponseWriter struct {
	Status int
	w      http.ResponseWriter
}

func newLoggedResponseWriter(w http.ResponseWriter) *loggedResponseWriter {
	lw := new(loggedResponseWriter)
	lw.w = w
	lw.Status = -1
	return lw
}

func (w *loggedResponseWriter) Flush() {
	if wf, ok := w.w.(http.Flusher); ok {
		wf.Flush()
	}
}

func (w *loggedResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *loggedResponseWriter) Write(d []byte) (int, error) {
	return w.w.Write(d)
}

func (w *loggedResponseWriter) WriteHeader(status int) {
	w.Status = status
	w.w.WriteHeader(status)
}
