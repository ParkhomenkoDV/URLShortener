package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestPostHandler(t *testing.T) {
	cfg := config.Config{BaseURL: "http://localhost:8080"}
	db := storage.New()
	h := New(&cfg, db)

	t.Run("plain text request", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader("https://example.com"))
		w := httptest.NewRecorder()
		h.Post(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		require.Contains(t, string(body), "http://localhost:8080/")
	})

	t.Run("JSON request", func(t *testing.T) {
		jsonBody := `{"url":"https://example.org"}`
		req := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.PostJSON(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)
		require.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var result struct {
			Result string `json:"result"`
		}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		require.Contains(t, result.Result, "http://localhost:8080/")
	})

	t.Run("invalid requests", func(t *testing.T) {
		tests := []struct {
			name        string
			body        string
			contentType string
			wantStatus  int
		}{
			{"empty body", "", "text/plain", http.StatusBadRequest},
			{"invalid URL", "not-a-url", "text/plain", http.StatusBadRequest},
			{"invalid JSON", "{bad}", "application/json", http.StatusBadRequest},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(tt.body))
				req.Header.Set("Content-Type", tt.contentType)
				w := httptest.NewRecorder()
				h.PostJSON(w, req)

				resp := w.Result()
				defer resp.Body.Close()
				require.Equal(t, tt.wantStatus, resp.StatusCode)
			})
		}
	})
}

func Test_normalizationURL(t *testing.T) {
	tests := []struct {
		name   string
		rawURL string
		want   string
	}{
		{"1", "yandex.ru", "http://yandex.ru"},
		{"2", "http://yandex.ru", "http://yandex.ru"},
		{"3", "https://yandex.ru", "https://yandex.ru"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizationURL(tt.rawURL)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_validateURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr bool
	}{
		{"1", "yandex.ru", false},
		{"2", "http://yandex.ru", false},
		{"3", "https://yandex.ru", false},
		{"4", "y a n d e x", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateURL(tt.rawURL); (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
