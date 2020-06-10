package httpapi

import (
	"net/http"
	"strings"
)

type OAuthChecker interface {
	Check(token string) bool
}

type OAuthMiddleware struct {
	Checker OAuthChecker;
}

func (middleware *OAuthMiddleware) Wrap(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		token := request.Header.Get("Authorization")
		if !strings.HasPrefix(token, "OAuth ") {
			writer.WriteHeader(http.StatusForbidden)
			return
		}

		token = token[len("OAuth "):]
		if !middleware.Checker.Check(token) {
			writer.WriteHeader(http.StatusForbidden)
			return
		}

		f(writer, request)
	}
}