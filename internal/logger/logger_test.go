package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggingMiddleware(t *testing.T) {
	t.Run("logging middleware", func(t *testing.T) {
		New()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("response"))
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:8080"

		rr := httptest.NewRecorder()
		middleware := LoggingMiddleware(handler)
		middleware.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "response", rr.Body.String())
	})

	t.Run("get client IP from headers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Real-IP", "192.168.1.1")

		ip := getClientIP(req)
		require.Equal(t, "192.168.1.1", ip)

		req.Header.Del("X-Real-IP")
		req.Header.Set("X-Forwarded-For", "10.0.0.1")

		ip = getClientIP(req)
		require.Equal(t, "10.0.0.1", ip)
	})
}
