package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (gz gzipResponseWriter) Write(data []byte) (int, error) {
	return gz.Writer.Write(data)
}

func GzipRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gzRead, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Unable to decode gzip body", http.StatusBadRequest)
				return
			}
			defer gzRead.Close()
			r.Body = gzRead
		}
		fmt.Println("GzipRequestMiddleware")
		next.ServeHTTP(w, r)
	})
}

func GzipResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gzWriter := gzip.NewWriter(w)
		defer gzWriter.Close()

		gz := gzipResponseWriter{Writer: gzWriter, ResponseWriter: w}
		fmt.Println("GzipResponseMiddleware")
		next.ServeHTTP(gz, r)
	})
}
