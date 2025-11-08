package middleware

import (
	"compress/gzip"

	"github.com/gin-gonic/gin"
)

// gzipWriter оборачивает gin.ResponseWriter для прозрачного сжатия исходящих данных в формате gzip.
// Перехватывает вызовы Write и WriteString, сжимая данные перед их записью в нижележащий ResponseWriter.
type gzipWriter struct {
	gin.ResponseWriter              // Встроенный ResponseWriter для доступа к базовым методам
	writer             *gzip.Writer // Gzip writer для сжатия данных
}

// Write сжимает переданные данные в формате gzip перед записью в нижележащий ResponseWriter.
// Возвращает количество записанных байт после сжатия и ошибку (если есть).
func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

// WriteString сжимает переданную строку в формате gzip перед записью в нижележащий ResponseWriter.
// Возвращает количество записанных байт после сжатия и ошибку (если есть).
func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

// Close закрывает gzip writer, гарантируя запись всех буферизованных данных.
// Должен вызываться для корректного завершения сжатия и избежания утечек ресурсов.
func (g *gzipWriter) Close() {
	g.writer.Close()
}
