package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type SimpleChecer struct {
}

func (checker SimpleChecer) Check(token string) bool {
	return token == "good"
}

func SimpleHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("ok"))
}

func TestOAuthMiddware(t *testing.T) {
	checker := SimpleChecer{}
	middleware := OAuthMiddleware{checker}

	handler := middleware.Wrap(SimpleHandler)
	t.Run("good", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		request.Header.Add("Authorization", "OAuth good")
		writer := httptest.NewRecorder()

		handler(writer, request)

		if writer.Code != http.StatusOK {
			t.Error("wrong code")
		}
	})

	t.Run("bad", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		request.Header.Add("Authorization", "OAuth bad")
		writer := httptest.NewRecorder()

		handler(writer, request)

		if writer.Code != http.StatusForbidden {
			t.Error("wrong code")
		}
	})

	t.Run("without", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		writer := httptest.NewRecorder()

		handler(writer, request)

		if writer.Code != http.StatusForbidden {
			t.Error("wrong code")
		}
	})
}
