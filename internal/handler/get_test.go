package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ParkhomenkoDV/URLShortener/internal/service/server"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {
	cfg := server.Config{BaseURL: "http://localhost:8080"}
	h := New(cfg)

	h.data["validID"] = "https://example.com"
	h.data["noScheme"] = "example.com"
	h.data["invalidURL"] = "http://invalid url.com"

	tests := []struct {
		name       string
		id         string
		setup      func(*Handler)
		wantStatus int
		wantURL    string
	}{
		{
			name:       "successful redirect",
			id:         "validID",
			wantStatus: http.StatusTemporaryRedirect,
			wantURL:    "https://example.com",
		}, {
			name:       "adds http scheme if missing",
			id:         "noScheme",
			wantStatus: http.StatusTemporaryRedirect,
			wantURL:    "http://example.com",
		}, {
			name:       "empty ID",
			id:         "",
			wantStatus: http.StatusBadRequest,
		}, {
			name:       "ID not found",
			id:         "nonExistentID",
			wantStatus: http.StatusNotFound,
		}, {
			name:       "invalid URL format",
			id:         "invalidURL",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/"+tt.id, nil)
			w := httptest.NewRecorder()

			// Создаем контекст с параметром маршрута
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			h.Get(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			require.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantURL != "" {
				require.Equal(t, tt.wantURL, resp.Header.Get("Location"))
			}
		})
	}
}
