package httputil

import (
	"bytes"
	"net/http"
	"strconv"
)

// Error will eventually be a shorthand for http.Error
// that takes just the http status code and creates
// a json response containing the verbose (yet uninformative)
// error desription. Currently plain string will do.
func Error(w http.ResponseWriter, status int) {
	var b bytes.Buffer
	b.WriteString(strconv.Itoa(status))
	b.WriteString(" ")
	b.WriteString(http.StatusText(status))
	http.Error(w, b.String(), status)
}
