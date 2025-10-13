package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware_Compression(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Инициализируем логгер для тестов
	InitLogger()

	router := gin.New()
	router.Use(GzipMiddleware())

	// Тестовый эндпоинт
	router.GET("/test", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	// Создаем запрос с поддержкой gzip
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем заголовки
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
	assert.Equal(t, "Accept-Encoding", w.Header().Get("Vary"))

	// Проверяем, что ответ сжат
	reader, err := gzip.NewReader(w.Body)
	assert.NoError(t, err)
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	assert.NoError(t, err)

	expectedJSON := `{"message":"Hello, World!"}`
	assert.JSONEq(t, expectedJSON, string(decompressed))
}

func TestGzipMiddleware_Decompression(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(GzipMiddleware())

	// Тестовый эндпоинт для POST
	router.POST("/test", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.String(http.StatusOK, string(body))
	})

	// Создаем сжатое тело запроса
	testData := "Hello, World!"
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Write([]byte(testData))
	writer.Close()

	// Создаем запрос со сжатым телом
	req := httptest.NewRequest("POST", "/test", &buf)
	req.Header.Set("Content-Encoding", "gzip")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем, что тело было правильно распаковано
	assert.Equal(t, testData, w.Body.String())
}

func TestGzipMiddleware_NoGzipSupport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(GzipMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Создаем запрос без поддержки gzip
	req := httptest.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем, что gzip не используется
	assert.Empty(t, w.Header().Get("Content-Encoding"))
	assert.Equal(t, "Hello, World!", w.Body.String())
}

func TestLoggingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(LoggingMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}
