package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	t.Run("Request struct creation and JSON marshaling", func(t *testing.T) {
		req := Request{
			URL: "https://example.com",
		}

		assert.Equal(t, "https://example.com", req.URL)

		// Тест JSON маршалинга
		jsonData, err := json.Marshal(req)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"url":"https://example.com"}`, string(jsonData))
	})

	t.Run("Request struct JSON unmarshaling", func(t *testing.T) {
		jsonStr := `{"url":"https://google.com"}`
		var req Request

		err := json.Unmarshal([]byte(jsonStr), &req)
		assert.NoError(t, err)
		assert.Equal(t, "https://google.com", req.URL)
	})
}

func TestResponse(t *testing.T) {
	t.Run("Response struct creation and JSON marshaling", func(t *testing.T) {
		resp := Response{
			Result: "https://short.ly/abc123",
		}

		assert.Equal(t, "https://short.ly/abc123", resp.Result)

		// Тест JSON маршалинга
		jsonData, err := json.Marshal(resp)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"result":"https://short.ly/abc123"}`, string(jsonData))
	})

	t.Run("Response struct JSON unmarshaling", func(t *testing.T) {
		jsonStr := `{"result":"https://short.ly/xyz789"}`
		var resp Response

		err := json.Unmarshal([]byte(jsonStr), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "https://short.ly/xyz789", resp.Result)
	})
}

func TestURLRecord(t *testing.T) {
	t.Run("URLRecord struct creation and JSON marshaling", func(t *testing.T) {
		record := URLRecord{
			ID:          1,
			ShortURL:    "abc123",
			OriginalURL: "https://example.com",
			UserID:      "user-123",
		}

		assert.Equal(t, 1, record.ID)
		assert.Equal(t, "abc123", record.ShortURL)
		assert.Equal(t, "https://example.com", record.OriginalURL)
		assert.Equal(t, "user-123", record.UserID)

		// Тест JSON маршалинга
		jsonData, err := json.Marshal(record)
		assert.NoError(t, err)
		expected := `{
			"id": 1,
			"short_url": "abc123",
			"original_url": "https://example.com",
			"user_id": "user-123"
		}`
		assert.JSONEq(t, expected, string(jsonData))
	})

	t.Run("URLRecord struct JSON unmarshaling", func(t *testing.T) {
		jsonStr := `{
			"id": 42,
			"short_url": "xyz789",
			"original_url": "https://google.com",
			"user_id": "user-456"
		}`
		var record URLRecord

		err := json.Unmarshal([]byte(jsonStr), &record)
		assert.NoError(t, err)
		assert.Equal(t, 42, record.ID)
		assert.Equal(t, "xyz789", record.ShortURL)
		assert.Equal(t, "https://google.com", record.OriginalURL)
		assert.Equal(t, "user-456", record.UserID)
	})
}

func TestBatchRequest(t *testing.T) {
	t.Run("BatchRequest struct creation and JSON marshaling", func(t *testing.T) {
		req := BatchRequest{
			CorrelationID: "req-1",
			OriginalURL:   "https://example.com",
		}

		assert.Equal(t, "req-1", req.CorrelationID)
		assert.Equal(t, "https://example.com", req.OriginalURL)

		// Тест JSON маршалинга
		jsonData, err := json.Marshal(req)
		assert.NoError(t, err)
		expected := `{
			"correlation_id": "req-1",
			"original_url": "https://example.com"
		}`
		assert.JSONEq(t, expected, string(jsonData))
	})

	t.Run("BatchRequest struct JSON unmarshaling", func(t *testing.T) {
		jsonStr := `{
			"correlation_id": "req-2",
			"original_url": "https://google.com"
		}`
		var req BatchRequest

		err := json.Unmarshal([]byte(jsonStr), &req)
		assert.NoError(t, err)
		assert.Equal(t, "req-2", req.CorrelationID)
		assert.Equal(t, "https://google.com", req.OriginalURL)
	})
}

func TestBatchResponse(t *testing.T) {
	t.Run("BatchResponse struct creation and JSON marshaling", func(t *testing.T) {
		resp := BatchResponse{
			CorrelationID: "req-1",
			ShortURL:      "https://short.ly/abc123",
		}

		assert.Equal(t, "req-1", resp.CorrelationID)
		assert.Equal(t, "https://short.ly/abc123", resp.ShortURL)

		// Тест JSON маршалинга
		jsonData, err := json.Marshal(resp)
		assert.NoError(t, err)
		expected := `{
			"correlation_id": "req-1",
			"short_url": "https://short.ly/abc123"
		}`
		assert.JSONEq(t, expected, string(jsonData))
	})

	t.Run("BatchResponse struct JSON unmarshaling", func(t *testing.T) {
		jsonStr := `{
			"correlation_id": "req-2",
			"short_url": "https://short.ly/xyz789"
		}`
		var resp BatchResponse

		err := json.Unmarshal([]byte(jsonStr), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "req-2", resp.CorrelationID)
		assert.Equal(t, "https://short.ly/xyz789", resp.ShortURL)
	})
}

func TestUserURL(t *testing.T) {
	t.Run("UserURL struct creation and JSON marshaling", func(t *testing.T) {
		userURL := UserURL{
			ShortURL:    "https://short.ly/abc123",
			OriginalURL: "https://example.com",
		}

		assert.Equal(t, "https://short.ly/abc123", userURL.ShortURL)
		assert.Equal(t, "https://example.com", userURL.OriginalURL)

		// Тест JSON маршалинга
		jsonData, err := json.Marshal(userURL)
		assert.NoError(t, err)
		expected := `{
			"short_url": "https://short.ly/abc123",
			"original_url": "https://example.com"
		}`
		assert.JSONEq(t, expected, string(jsonData))
	})

	t.Run("UserURL struct JSON unmarshaling", func(t *testing.T) {
		jsonStr := `{
			"short_url": "https://short.ly/xyz789",
			"original_url": "https://google.com"
		}`
		var userURL UserURL

		err := json.Unmarshal([]byte(jsonStr), &userURL)
		assert.NoError(t, err)
		assert.Equal(t, "https://short.ly/xyz789", userURL.ShortURL)
		assert.Equal(t, "https://google.com", userURL.OriginalURL)
	})
}
