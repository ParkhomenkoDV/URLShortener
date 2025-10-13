package middleware

import (
	"compress/gzip"

	"github.com/gin-gonic/gin"
)

// gzipWriter оборачивает ResponseWriter для сжатия gzip
type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

// Write сжимает данные перед записью
func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

// WriteString сжимает строку перед записью
func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

// Close закрывает gzip writer
func (g *gzipWriter) Close() {
	g.writer.Close()
}