package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/handler"
	"github.com/ParkhomenkoDV/URLShortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	cfg := config.Config{BaseURL: "http://localhost:8080"}
	db := storage.New()
	h := handler.New(&cfg, db)

	t.Run("full flow: post then get", func(t *testing.T) {
		// Post request to shorten URL
		postReq := httptest.NewRequest("POST", "/", strings.NewReader("https://example.com"))
		postRec := httptest.NewRecorder()
		h.Post(postRec, postReq)

		require.Equal(t, http.StatusCreated, postRec.Code)
		shortURL := strings.TrimSpace(postRec.Body.String())

		// Extract ID from short URL
		id := strings.TrimPrefix(shortURL, cfg.BaseURL+"/")

		// Get request to redirect
		getReq := httptest.NewRequest("GET", "/"+id, nil)
		getRec := httptest.NewRecorder()

		// Use chi context for URL params
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", id)
		getReq = getReq.WithContext(context.WithValue(getReq.Context(), chi.RouteCtxKey, rctx))

		h.Get(getRec, getReq)

		require.Equal(t, http.StatusTemporaryRedirect, getRec.Code)
		require.Equal(t, "https://example.com", getRec.Header().Get("Location"))
	})

	t.Run("JSON API flow", func(t *testing.T) {
		jsonBody := `{"url":"https://google.com"}`
		postReq := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(jsonBody))
		postReq.Header.Set("Content-Type", "application/json")
		postRec := httptest.NewRecorder()
		h.PostJSON(postRec, postReq)

		require.Equal(t, http.StatusCreated, postRec.Code)

		var response struct {
			Result string `json:"result"`
		}
		err := json.NewDecoder(postRec.Body).Decode(&response)
		require.NoError(t, err)
		require.Contains(t, response.Result, cfg.BaseURL)
	})
}
