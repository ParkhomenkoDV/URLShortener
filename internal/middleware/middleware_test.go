package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGzipMiddleware(t *testing.T) {
	t.Run("gzip response compression", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("test response content"))
		})

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept-Encoding", "gzip")

		rec := httptest.NewRecorder()
		middleware := GzipResponseMiddleware(handler)
		middleware.ServeHTTP(rec, req)

		require.Equal(t, "gzip", rec.Header().Get("Content-Encoding"))
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("gzip request decompression", func(t *testing.T) {
		// Create gzipped content
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		gz.Write([]byte("compressed content"))
		gz.Close()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			require.Equal(t, "compressed content", string(body))
			w.Write([]byte("ok"))
		})

		req := httptest.NewRequest("POST", "/", &buf)
		req.Header.Set("Content-Encoding", "gzip")

		rec := httptest.NewRecorder()
		middleware := GzipRequestMiddleware(handler)
		middleware.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("no compression when not accepted", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("test response"))
		})

		req := httptest.NewRequest("GET", "/", nil)
		// No Accept-Encoding header

		rec := httptest.NewRecorder()
		middleware := GzipResponseMiddleware(handler)
		middleware.ServeHTTP(rec, req)

		require.Empty(t, rec.Header().Get("Content-Encoding"))
		require.Equal(t, "test response", rec.Body.String())
	})
}
