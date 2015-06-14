package middleware

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
	"github.com/vilppuvuorinen/chalice/httputil"
)

type key int

var authUsernameKey key = 0

// NewAnyAuth creates a new auth decorator accepting any of
// the supported methods for credential passing. Validation
// is performed with provided function.
func NewAnyAuth(
	authFunc func(username, password string) bool,
) func(
	func(context.Context, http.ResponseWriter, *http.Request),
) func(context.Context, http.ResponseWriter, *http.Request) {
	return newAuth(parseAny, authFunc)
}

// NewPostJsonAuth creates a new auth decorator accepting
// credentials passed in POST body in format:
//   { "Username": "user", "Password": "pass" }
// Validation is performed with provided function.
func NewPostJsonAuth(
	authFunc func(username, password string) bool,
) func(
	func(context.Context, http.ResponseWriter, *http.Request),
) func(context.Context, http.ResponseWriter, *http.Request) {
	return newAuth(parsePostJson, authFunc)
}

// NewBasicAuth creates a new auth decorator accepting
// BasicAuth credentials. Validation is performed with
// provided function.
func NewBasicAuth(
	authFunc func(username, password string) bool,
) func(
	func(context.Context, http.ResponseWriter, *http.Request),
) func(context.Context, http.ResponseWriter, *http.Request) {
	return newAuth(parseBasic, authFunc)
}

// GetUsername retrieves username from given context if one exists.
func GetUsername(c context.Context) (username string, ok bool) {
	username, ok = c.Value(authUsernameKey).(string)
	return username, ok
}

func newAuth(
	parseFunc func(r *http.Request) (username, password string, ok bool),
	authFunc func(username, password string) bool,
) func(
	func(context.Context, http.ResponseWriter, *http.Request),
) func(context.Context, http.ResponseWriter, *http.Request) {

	return func(
		f func(context.Context, http.ResponseWriter, *http.Request),
	) func(context.Context, http.ResponseWriter, *http.Request) {

		return func(c context.Context, w http.ResponseWriter, r *http.Request) {
			username, password, ok := parseFunc(r)
			if !ok {
				httputil.Error(w, http.StatusBadRequest)
				return
			}

			if authFunc(username, password) {
				c := context.WithValue(c, authUsernameKey, username)
				f(c, w, r)
			} else {
				httputil.Error(w, http.StatusUnauthorized)
			}
		}
	}
}

func parseBasic(r *http.Request) (username, password string, ok bool) {
	return r.BasicAuth()
}

type auth struct {
	Username, Password string
}

func parsePostJson(r *http.Request) (username, password string, ok bool) {
	d := json.NewDecoder(r.Body)
	var a auth
	err := d.Decode(&a)
	if err != nil {
		return "", "", false
	}
	return a.Username, a.Password, true
}

func parseAny(r *http.Request) (username, password string, ok bool) {
	u, p, k := parseBasic(r)
	if ok {
		return u, p, k
	}

	// No need to verify. Failure will be relayed.
	return parsePostJson(r)
}
