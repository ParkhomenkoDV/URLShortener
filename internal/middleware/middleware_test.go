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

// TestGzipMiddleware_Compression тестирует сжатие исходящих ответов в формате gzip.
// Проверяет что:
// - Устанавливаются правильные заголовки Content-Encoding и Vary
// - Данные правильно сжимаются и могут быть распакованы
// - Содержимое ответа сохраняется корректно после сжатия/распаковки
func TestGzipMiddleware_Compression(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Инициализируем логгер для тестов
	InitLogger()

	router := gin.New()
	router.Use(GzipMiddleware())

	// Тестовый endpoint возвращающий JSON данные
	router.GET("/test", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	// Создаем запрос с указанием поддержки gzip в Accept-Encoding
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ПРОВЕРКА ЗАГОЛОВКОВ ОТВЕТА
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"),
		"Content-Encoding должен быть 'gzip' для сжатого ответа")
	assert.Equal(t, "Accept-Encoding", w.Header().Get("Vary"),
		"Vary должен содержать 'Accept-Encoding' для кэширующих прокси")

	// ПРОВЕРКА СЖАТЫХ ДАННЫХ
	// Создаем reader для распаковки gzip данных
	reader, err := gzip.NewReader(w.Body)
	assert.NoError(t, err, "Должен создаваться gzip reader для сжатого ответа")
	defer reader.Close()

	// Читаем и распаковываем данные
	decompressed, err := io.ReadAll(reader)
	assert.NoError(t, err, "Должны успешно читаться распакованные данные")

	// Проверяем что распакованные данные соответствуют ожидаемому JSON
	expectedJSON := `{"message":"Hello, World!"}`
	assert.JSONEq(t, expectedJSON, string(decompressed),
		"Распакованные данные должны соответствовать исходному JSON")
}

// TestGzipMiddleware_Decompression тестирует распаковку входящих сжатых запросов.
// Проверяет что middleware корректно распаковывает gzip тело запроса перед передачей обработчику.
func TestGzipMiddleware_Decompression(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(GzipMiddleware())

	// Тестовый endpoint для POST запросов, который возвращает полученное тело
	router.POST("/test", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.String(http.StatusOK, string(body))
	})

	// ПОДГОТОВКА СЖАТЫХ ДАННЫХ ДЛЯ ЗАПРОСА
	testData := "Hello, World!"
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Write([]byte(testData))
	writer.Close()

	// Создаем запрос со сжатым телом и указанием gzip в Content-Encoding
	req := httptest.NewRequest("POST", "/test", &buf)
	req.Header.Set("Content-Encoding", "gzip")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем, что middleware корректно распаковал тело запроса
	assert.Equal(t, testData, w.Body.String(),
		"Обработчик должен получить распакованные исходные данные")
}

// TestGzipMiddleware_NoGzipSupport тестирует поведение middleware когда клиент не поддерживает gzip.
// Проверяет что:
// - Ответ не сжимается когда клиент не указывает Accept-Encoding: gzip
// - Данные передаются в исходном виде
// - Заголовки Content-Encoding не устанавливаются
func TestGzipMiddleware_NoGzipSupport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(GzipMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Создаем запрос БЕЗ поддержки gzip
	req := httptest.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ПРОВЕРКА ОТСУТСТВИЯ СЖАТИЯ
	assert.Empty(t, w.Header().Get("Content-Encoding"),
		"Content-Encoding должен быть пустым когда клиент не поддерживает gzip")
	assert.Equal(t, "Hello, World!", w.Body.String(),
		"Тело ответа должно передаваться без сжатия")
}

// TestLoggingMiddleware тестирует middleware логирования.
// Проверяет что запросы корректно обрабатываются и логируются без изменения функциональности.
func TestLoggingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(LoggingMiddleware())

	// Простой endpoint для тестирования
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Проверяем что middleware не влияет на нормальную обработку запросов
	assert.Equal(t, http.StatusOK, w.Code, "Должен возвращаться статус 200 OK")
	assert.Equal(t, "OK", w.Body.String(), "Тело ответа должно быть корректным")
}
