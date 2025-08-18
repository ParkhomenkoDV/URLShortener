package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ParkhomenkoDV/URLShortener/internal/utils"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalURL := strings.TrimSpace(string(body))
	if err := validateURL(originalURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortKey := utils.GenerateShortURL(8)
	shortURL := h.config.BaseURL + "/" + shortKey

	h.mutex.Lock()
	h.data[shortKey] = originalURL
	h.mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *Handler) PostJSON(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := validateURL(req.URL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortKey := utils.GenerateShortURL(8)
	shortURL := h.config.BaseURL + "/" + shortKey

	h.mutex.Lock()
	h.data[shortKey] = req.URL
	h.mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ShortenResponse{Result: shortURL})
}

// normalizationURL - нормализация url.
func normalizationURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)

	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "http://" + rawURL
	}

	return rawURL
}

// validateURL - валидация url.
func validateURL(rawURL string) error {
	rawURL = normalizationURL(rawURL)

	if _, err := url.ParseRequestURI(rawURL); err != nil {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}
