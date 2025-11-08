package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/repository"
	"github.com/ParkhomenkoDV/URLShortener/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTest создает тестовый HTTP-сервер с инициализированными зависимостями.
// Возвращает роутер Gin и обработчик для прямого вызова методов сервиса.
func setupTest() (*gin.Engine, *Handler) {
	// Инициализируем in-memory репозиторий для изоляции тестов
	repo := repository.NewMemory()

	// Создаем тестовую конфигурацию
	configuration := &config.Config{
		Port:         ":8080",
		ShortAddress: "http://localhost:8080",
		LengthKey:    6,
	}

	// Инициализируем сервисный слой
	service := service.New(repo, configuration)

	// Создаем роутер Gin
	ginEngine := gin.Default()

	// Создаем обработчик (для прямого доступа в тестах)
	h := &Handler{Service: service}

	// Инициализируем маршруты и middleware
	New(ginEngine, service, configuration)

	return ginEngine, h
}

// TestSendURLHandler тестирует обработчик создания коротких ссылок через text/plain endpoint
func TestSendURLHandler(t *testing.T) {
	mux, _ := setupTest()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Тест 1: Успешное создание короткой ссылки
	t.Run("successful short URL creation with text/plain", func(t *testing.T) {
		longURL := "https://example.com/very/long/url"
		req, err := http.NewRequest("POST", server.URL+"/", bytes.NewBufferString(longURL))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Должен вернуться статус 201 Created")
		assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"), "Content-Type должен быть text/plain")

		// Читаем ответ, но не проверяем конкретное значение (оно генерируется случайно)
		_, err = io.ReadAll(resp.Body)
		assert.NoError(t, err)
	})

	// Тест 2: Неверный Content-Type (должен быть text/plain)
	t.Run("reject request with invalid content type", func(t *testing.T) {
		req, err := http.NewRequest("POST", server.URL+"/", bytes.NewBufferString("https://test.com"))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json") // Неправильный Content-Type

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Должен вернуться статус 400 Bad Request")
	})

	// Тест 3: Пустое тело запроса (граничный случай)
	t.Run("handle empty request body", func(t *testing.T) {
		req, err := http.NewRequest("POST", server.URL+"/", bytes.NewBufferString(""))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Ожидаем 201 даже для пустого URL (поведение зависит от реализации сервиса)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})
}

// TestGetURLHandler тестирует обработчик редиректа по коротким ссылкам
func TestGetURLHandler(t *testing.T) {
	mux, h := setupTest()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Предварительно создаем тестовую короткую ссылку через сервис
	longURL := "https://redirect.me"
	shortURL, err := h.Service.CreateShortURL(longURL, "")
	assert.NoError(t, err)

	// Создаем HTTP-клиент, который не следует за редиректами автоматически
	// Это позволяет проверить код статуса и заголовок Location
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Останавливаемся после первого редиректа
		},
	}

	// Тест 1: Успешный редирект по существующей короткой ссылке
	t.Run("successful redirect for existing short URL", func(t *testing.T) {
		req, err := http.NewRequest("GET", server.URL+"/"+shortURL, nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode, "Должен вернуться статус 307 Temporary Redirect")
		assert.Equal(t, longURL, resp.Header.Get("Location"), "Заголовок Location должен содержать оригинальный URL")
	})

	// Тест 2: Запрос несуществующей короткой ссылки
	t.Run("return error for non-existent short URL", func(t *testing.T) {
		req, err := http.NewRequest("GET", server.URL+"/nonexistent", nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Должен вернуться статус 400 Bad Request для несуществующей ссылки")
	})
}

// TestSendJSONURLHandler тестирует обработчик создания коротких ссылок через JSON endpoint
func TestSendJSONURLHandler(t *testing.T) {
	mux, _ := setupTest()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Тест 1: Успешное создание короткой ссылки через JSON API
	t.Run("successful short URL creation via JSON API", func(t *testing.T) {
		jsonBody := `{"url": "https://example.com/very/long/url/json"}`
		req, err := http.NewRequest("POST", server.URL+"/api/shorten", bytes.NewBufferString(jsonBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Должен вернуться статус 201 Created")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type должен быть application/json")

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		// Проверяем структуру JSON ответа
		assert.Contains(t, string(body), `"result"`, "Ответ должен содержать поле result")
		assert.Contains(t, string(body), "http://localhost:8080/", "Ответ должен содержать базовый URL")
	})

	// Тест 2: Неверный Content-Type для JSON endpoint'а
	t.Run("reject JSON request with invalid content type", func(t *testing.T) {
		jsonBody := `{"url": "https://test.com"}`
		req, err := http.NewRequest("POST", server.URL+"/api/shorten", bytes.NewBufferString(jsonBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain") // Неправильный Content-Type

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Должен вернуться статус 400 Bad Request")
	})

	// Тест 3: Невалидный JSON в теле запроса
	t.Run("reject request with invalid JSON syntax", func(t *testing.T) {
		invalidJSON := `{"url": "https://test.com",}` // trailing comma - синтаксическая ошибка
		req, err := http.NewRequest("POST", server.URL+"/api/shorten", bytes.NewBufferString(invalidJSON))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Должен вернуться статус 400 Bad Request для невалидного JSON")
	})

	// Тест 4: Проверка формата JSON ответа
	t.Run("validate JSON response format", func(t *testing.T) {
		jsonBody := `{"url": "https://google.com"}`
		req, err := http.NewRequest("POST", server.URL+"/api/shorten", bytes.NewBufferString(jsonBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response struct {
			Result string `json:"result"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// Проверяем корректность формата ответа
		assert.NotEmpty(t, response.Result, "Поле result не должно быть пустым")
		assert.True(t, strings.HasPrefix(response.Result, "http://localhost:8080/"),
			"Ответ должен начинаться с базового URL")
	})
}

// TestBothEndpointsCreateURLs интеграционный тест для обоих endpoint'ов создания ссылок
func TestBothEndpointsCreateURLs(t *testing.T) {
	mux, _ := setupTest()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Тест 1: Текстовый endpoint создает короткую ссылку
	t.Run("text endpoint creates short URL successfully", func(t *testing.T) {
		longURL := "https://text-endpoint-test.com"
		req, err := http.NewRequest("POST", server.URL+"/", bytes.NewBufferString(longURL))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		shortURLBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		shortURL := strings.TrimSpace(string(shortURLBytes))

		// Проверяем что что-то вернулось (конкретное значение зависит от реализации)
		assert.NotEmpty(t, shortURL, "Должен вернуться непустой короткий URL")
	})

	// Тест 2: JSON endpoint создает короткую ссылку
	t.Run("JSON endpoint creates short URL successfully", func(t *testing.T) {
		longURL := "https://json-endpoint-test.com"
		jsonBody := `{"url": "` + longURL + `"}`
		req, err := http.NewRequest("POST", server.URL+"/api/shorten", bytes.NewBufferString(jsonBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response struct {
			Result string `json:"result"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// Проверяем формат ответа
		assert.NotEmpty(t, response.Result, "Поле result не должно быть пустым")
		assert.True(t, strings.HasPrefix(response.Result, "http://localhost:8080/"),
			"Ответ должен содержать базовый URL")
	})
}

// TestRedirectIntegration тестирует интеграцию создания ссылки и редиректа
func TestRedirectIntegration(t *testing.T) {
	mux, h := setupTest()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Создаем клиент без авто-редиректов для проверки промежуточных ответов
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	t.Run("complete flow: create short URL and redirect", func(t *testing.T) {
		// Создаем ссылку напрямую через сервис (для контроля)
		longURL := "https://redirect-test.com"
		shortURL, err := h.Service.CreateShortURL(longURL, "")
		assert.NoError(t, err)

		// Проверяем редирект через HTTP endpoint
		req, err := http.NewRequest("GET", server.URL+"/"+shortURL, nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode, "Должен вернуться статус 307 Temporary Redirect")
		assert.Equal(t, longURL, resp.Header.Get("Location"), "Location должен соответствовать оригинальному URL")
	})
}

// TestSendJSONURLBatchHandler тестирует обработчик пакетного создания коротких ссылок
func TestSendJSONURLBatchHandler(t *testing.T) {
	mux, _ := setupTest()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Тест 1: Успешное пакетное создание нескольких коротких ссылок
	t.Run("successful batch creation of multiple short URLs", func(t *testing.T) {
		jsonBody := `[
			{"correlation_id": "1", "original_url": "https://example1.com"},
			{"correlation_id": "2", "original_url": "https://example2.com"},
			{"correlation_id": "3", "original_url": "https://example3.com"}
		]`
		req, err := http.NewRequest("POST", server.URL+"/api/shorten/batch", bytes.NewBufferString(jsonBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Должен вернуться статус 201 Created")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type должен быть application/json")

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var responses []map[string]string
		err = json.Unmarshal(body, &responses)
		assert.NoError(t, err)

		// Проверяем количество ответов
		assert.Len(t, responses, 3, "Должно вернуться 3 ответа")

		// Проверяем что все correlation_id присутствуют в ответах
		correlationIDs := map[string]bool{}
		for _, response := range responses {
			correlationIDs[response["correlation_id"]] = true
			assert.Contains(t, response["short_url"], "http://localhost:8080/",
				"Каждый короткий URL должен содержать базовый адрес")
		}
		assert.True(t, correlationIDs["1"], "Должен присутствовать correlation_id '1'")
		assert.True(t, correlationIDs["2"], "Должен присутствовать correlation_id '2'")
		assert.True(t, correlationIDs["3"], "Должен присутствовать correlation_id '3'")
	})

	// Тест 2: Пустой пакетный запрос
	t.Run("reject empty batch request", func(t *testing.T) {
		jsonBody := `[]`
		req, err := http.NewRequest("POST", server.URL+"/api/shorten/batch", bytes.NewBufferString(jsonBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Должен вернуться статус 400 Bad Request для пустого пакета")
	})

	// Тест 3: Большой пакетный запрос (проверка производительности)
	t.Run("handle large batch request", func(t *testing.T) {
		var requests []map[string]string
		for i := 1; i <= 100; i++ {
			requests = append(requests, map[string]string{
				"correlation_id": fmt.Sprintf("id_%d", i),
				"original_url":   fmt.Sprintf("https://example%d.com", i),
			})
		}

		jsonBody, err := json.Marshal(requests)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", server.URL+"/api/shorten/batch", bytes.NewBufferString(string(jsonBody)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Должен вернуться статус 201 Created для большого пакета")

		var responses []map[string]string
		err = json.NewDecoder(resp.Body).Decode(&responses)
		assert.NoError(t, err)

		// Проверяем что все 100 запросов обработаны
		assert.Len(t, responses, 100, "Должно вернуться 100 ответов")
	})
}
