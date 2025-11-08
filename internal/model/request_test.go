package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRequest тестирует структуру Request и её JSON маршалинг/анмаршалинг
func TestRequest(t *testing.T) {
	// Тест 1: Создание структуры Request и проверка полей
	t.Run("Request struct creation and field assignment", func(t *testing.T) {
		req := Request{
			URL: "https://example.com",
		}

		assert.Equal(t, "https://example.com", req.URL,
			"Поле URL должно содержать переданное значение")
	})

	// Тест 2: JSON маршалинг структуры Request
	t.Run("Request JSON marshaling produces correct JSON", func(t *testing.T) {
		req := Request{
			URL: "https://example.com",
		}

		jsonData, err := json.Marshal(req)
		assert.NoError(t, err, "Маршалинг Request в JSON не должен возвращать ошибку")
		assert.JSONEq(t, `{"url":"https://example.com"}`, string(jsonData),
			"Сгенерированный JSON должен соответствовать ожидаемому формату")
	})

	// Тест 3: JSON анмаршалинг в структуру Request
	t.Run("Request JSON unmarshaling correctly populates struct", func(t *testing.T) {
		jsonStr := `{"url":"https://google.com"}`
		var req Request

		err := json.Unmarshal([]byte(jsonStr), &req)
		assert.NoError(t, err, "Анмаршалинг JSON в Request не должен возвращать ошибку")
		assert.Equal(t, "https://google.com", req.URL,
			"Поле URL должно быть корректно заполнено из JSON")
	})
}

// TestResponse тестирует структуру Response и её JSON маршалинг/анмаршалинг
func TestResponse(t *testing.T) {
	// Тест 1: Создание структуры Response и проверка полей
	t.Run("Response struct creation and field assignment", func(t *testing.T) {
		resp := Response{
			Result: "https://short.ly/abc123",
		}

		assert.Equal(t, "https://short.ly/abc123", resp.Result,
			"Поле Result должно содержать переданное значение")
	})

	// Тест 2: JSON маршалинг структуры Response
	t.Run("Response JSON marshaling produces correct JSON", func(t *testing.T) {
		resp := Response{
			Result: "https://short.ly/abc123",
		}

		jsonData, err := json.Marshal(resp)
		assert.NoError(t, err, "Маршалинг Response в JSON не должен возвращать ошибку")
		assert.JSONEq(t, `{"result":"https://short.ly/abc123"}`, string(jsonData),
			"Сгенерированный JSON должен соответствовать ожидаемому формату")
	})

	// Тест 3: JSON анмаршалинг в структуру Response
	t.Run("Response JSON unmarshaling correctly populates struct", func(t *testing.T) {
		jsonStr := `{"result":"https://short.ly/xyz789"}`
		var resp Response

		err := json.Unmarshal([]byte(jsonStr), &resp)
		assert.NoError(t, err, "Анмаршалинг JSON в Response не должен возвращать ошибку")
		assert.Equal(t, "https://short.ly/xyz789", resp.Result,
			"Поле Result должно быть корректно заполнено из JSON")
	})
}

// TestURLRecord тестирует структуру URLRecord и её JSON маршалинг/анмаршалинг
func TestURLRecord(t *testing.T) {
	// Тест 1: Создание структуры URLRecord и проверка всех полей
	t.Run("URLRecord struct creation and all fields assignment", func(t *testing.T) {
		record := URLRecord{
			ID:          1,
			ShortURL:    "abc123",
			OriginalURL: "https://example.com",
			UserID:      "user-123",
		}

		assert.Equal(t, 1, record.ID, "Поле ID должно содержать переданное значение")
		assert.Equal(t, "abc123", record.ShortURL, "Поле ShortURL должно содержать переданное значение")
		assert.Equal(t, "https://example.com", record.OriginalURL, "Поле OriginalURL должно содержать переданное значение")
		assert.Equal(t, "user-123", record.UserID, "Поле UserID должно содержать переданное значение")
	})

	// Тест 2: JSON маршалинг структуры URLRecord
	t.Run("URLRecord JSON marshaling produces correct JSON with all fields", func(t *testing.T) {
		record := URLRecord{
			ID:          1,
			ShortURL:    "abc123",
			OriginalURL: "https://example.com",
			UserID:      "user-123",
		}

		jsonData, err := json.Marshal(record)
		assert.NoError(t, err, "Маршалинг URLRecord в JSON не должен возвращать ошибку")

		expected := `{
			"id": 1,
			"short_url": "abc123",
			"original_url": "https://example.com",
			"user_id": "user-123"
		}`
		assert.JSONEq(t, expected, string(jsonData),
			"Сгенерированный JSON должен содержать все поля в правильном формате")
	})

	// Тест 3: JSON анмаршалинг в структуру URLRecord
	t.Run("URLRecord JSON unmarshaling correctly populates all fields", func(t *testing.T) {
		jsonStr := `{
			"id": 42,
			"short_url": "xyz789",
			"original_url": "https://google.com",
			"user_id": "user-456"
		}`
		var record URLRecord

		err := json.Unmarshal([]byte(jsonStr), &record)
		assert.NoError(t, err, "Анмаршалинг JSON в URLRecord не должен возвращать ошибку")
		assert.Equal(t, 42, record.ID, "Поле ID должно быть корректно заполнено из JSON")
		assert.Equal(t, "xyz789", record.ShortURL, "Поле ShortURL должно быть корректно заполнено из JSON")
		assert.Equal(t, "https://google.com", record.OriginalURL, "Поле OriginalURL должно быть корректно заполнено из JSON")
		assert.Equal(t, "user-456", record.UserID, "Поле UserID должно быть корректно заполнено из JSON")
	})
}

// TestBatchRequest тестирует структуру BatchRequest и её JSON маршалинг/анмаршалинг
func TestBatchRequest(t *testing.T) {
	// Тест 1: Создание структуры BatchRequest и проверка полей
	t.Run("BatchRequest struct creation and field assignment", func(t *testing.T) {
		req := BatchRequest{
			CorrelationID: "req-1",
			OriginalURL:   "https://example.com",
		}

		assert.Equal(t, "req-1", req.CorrelationID,
			"Поле CorrelationID должно содержать переданное значение")
		assert.Equal(t, "https://example.com", req.OriginalURL,
			"Поле OriginalURL должно содержать переданное значение")
	})

	// Тест 2: JSON маршалинг структуры BatchRequest
	t.Run("BatchRequest JSON marshaling produces correct JSON", func(t *testing.T) {
		req := BatchRequest{
			CorrelationID: "req-1",
			OriginalURL:   "https://example.com",
		}

		jsonData, err := json.Marshal(req)
		assert.NoError(t, err, "Маршалинг BatchRequest в JSON не должен возвращать ошибку")

		expected := `{
			"correlation_id": "req-1",
			"original_url": "https://example.com"
		}`
		assert.JSONEq(t, expected, string(jsonData),
			"Сгенерированный JSON должен соответствовать ожидаемому формату")
	})

	// Тест 3: JSON анмаршалинг в структуру BatchRequest
	t.Run("BatchRequest JSON unmarshaling correctly populates struct", func(t *testing.T) {
		jsonStr := `{
			"correlation_id": "req-2",
			"original_url": "https://google.com"
		}`
		var req BatchRequest

		err := json.Unmarshal([]byte(jsonStr), &req)
		assert.NoError(t, err, "Анмаршалинг JSON в BatchRequest не должен возвращать ошибку")
		assert.Equal(t, "req-2", req.CorrelationID,
			"Поле CorrelationID должно быть корректно заполнено из JSON")
		assert.Equal(t, "https://google.com", req.OriginalURL,
			"Поле OriginalURL должно быть корректно заполнено из JSON")
	})
}

// TestBatchResponse тестирует структуру BatchResponse и её JSON маршалинг/анмаршалинг
func TestBatchResponse(t *testing.T) {
	// Тест 1: Создание структуры BatchResponse и проверка полей
	t.Run("BatchResponse struct creation and field assignment", func(t *testing.T) {
		resp := BatchResponse{
			CorrelationID: "req-1",
			ShortURL:      "https://short.ly/abc123",
		}

		assert.Equal(t, "req-1", resp.CorrelationID,
			"Поле CorrelationID должно содержать переданное значение")
		assert.Equal(t, "https://short.ly/abc123", resp.ShortURL,
			"Поле ShortURL должно содержать переданное значение")
	})

	// Тест 2: JSON маршалинг структуры BatchResponse
	t.Run("BatchResponse JSON marshaling produces correct JSON", func(t *testing.T) {
		resp := BatchResponse{
			CorrelationID: "req-1",
			ShortURL:      "https://short.ly/abc123",
		}

		jsonData, err := json.Marshal(resp)
		assert.NoError(t, err, "Маршалинг BatchResponse в JSON не должен возвращать ошибку")

		expected := `{
			"correlation_id": "req-1",
			"short_url": "https://short.ly/abc123"
		}`
		assert.JSONEq(t, expected, string(jsonData),
			"Сгенерированный JSON должен соответствовать ожидаемому формату")
	})

	// Тест 3: JSON анмаршалинг в структуру BatchResponse
	t.Run("BatchResponse JSON unmarshaling correctly populates struct", func(t *testing.T) {
		jsonStr := `{
			"correlation_id": "req-2",
			"short_url": "https://short.ly/xyz789"
		}`
		var resp BatchResponse

		err := json.Unmarshal([]byte(jsonStr), &resp)
		assert.NoError(t, err, "Анмаршалинг JSON в BatchResponse не должен возвращать ошибку")
		assert.Equal(t, "req-2", resp.CorrelationID,
			"Поле CorrelationID должно быть корректно заполнено из JSON")
		assert.Equal(t, "https://short.ly/xyz789", resp.ShortURL,
			"Поле ShortURL должно быть корректно заполнено из JSON")
	})
}

// TestUserURL тестирует структуру UserURL и её JSON маршалинг/анмаршалинг
func TestUserURL(t *testing.T) {
	// Тест 1: Создание структуры UserURL и проверка полей
	t.Run("UserURL struct creation and field assignment", func(t *testing.T) {
		userURL := UserURL{
			ShortURL:    "https://short.ly/abc123",
			OriginalURL: "https://example.com",
		}

		assert.Equal(t, "https://short.ly/abc123", userURL.ShortURL,
			"Поле ShortURL должно содержать переданное значение")
		assert.Equal(t, "https://example.com", userURL.OriginalURL,
			"Поле OriginalURL должно содержать переданное значение")
	})

	// Тест 2: JSON маршалинг структуры UserURL
	t.Run("UserURL JSON marshaling produces correct JSON", func(t *testing.T) {
		userURL := UserURL{
			ShortURL:    "https://short.ly/abc123",
			OriginalURL: "https://example.com",
		}

		jsonData, err := json.Marshal(userURL)
		assert.NoError(t, err, "Маршалинг UserURL в JSON не должен возвращать ошибку")

		expected := `{
			"short_url": "https://short.ly/abc123",
			"original_url": "https://example.com"
		}`
		assert.JSONEq(t, expected, string(jsonData),
			"Сгенерированный JSON должен соответствовать ожидаемому формату")
	})

	// Тест 3: JSON анмаршалинг в структуру UserURL
	t.Run("UserURL JSON unmarshaling correctly populates struct", func(t *testing.T) {
		jsonStr := `{
			"short_url": "https://short.ly/xyz789",
			"original_url": "https://google.com"
		}`
		var userURL UserURL

		err := json.Unmarshal([]byte(jsonStr), &userURL)
		assert.NoError(t, err, "Анмаршалинг JSON в UserURL не должен возвращать ошибку")
		assert.Equal(t, "https://short.ly/xyz789", userURL.ShortURL,
			"Поле ShortURL должно быть корректно заполнено из JSON")
		assert.Equal(t, "https://google.com", userURL.OriginalURL,
			"Поле OriginalURL должно быть корректно заполнено из JSON")
	})
}
