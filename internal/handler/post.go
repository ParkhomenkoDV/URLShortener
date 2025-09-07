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
	shortURL, err := h.processURL(originalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *Handler) PostJSON(w http.ResponseWriter, r *http.Request) {
	// Проверяем Content-Type
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	shortURL, err := h.processURL(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

func (h *Handler) processURL(rawURL string) (string, error) {
	if err := validateURL(rawURL); err != nil {
		return "", err
	}

	var shortKey string
	for {
		shortKey = utils.GenerateShortURL(8)
		h.mutex.Lock()
		_, exists := h.data[shortKey]
		h.mutex.Unlock()
		if !exists {
			break
		}
	}

	shortURL := h.config.BaseURL + "/" + shortKey

	h.mutex.Lock()
	h.data[shortKey] = normalizationURL(rawURL)
	h.mutex.Unlock()

	return shortURL, nil
}
