package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimSpace(r.Header.Get("Authorization"))
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		expectedToken := strings.TrimSpace(os.Getenv("WEB_STREAMING_AUTH_TOKEN"))
		if expectedToken == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		providedToken := strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
		if subtle.ConstantTimeCompare([]byte(providedToken), []byte(expectedToken)) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
