package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ParkhomenkoDV/URLShortener/internal/model"
	"github.com/ParkhomenkoDV/URLShortener/internal/utils"
)

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
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	var req model.Request
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
	json.NewEncoder(w).Encode(model.Response{Result: shortURL})
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

// processURL - сокращение + запись URL.
func (h *Handler) processURL(rawURL string) (string, error) {
	if err := validateURL(rawURL); err != nil {
		return "", err
	}

	var shortKey string
	for {
		shortKey = utils.GenerateShortURL(8)
		_, exists := h.db.Get(shortKey)
		if !exists {
			break
		}
	}

	shortURL := h.config.BaseURL + "/" + shortKey

	h.db.Set(shortKey, normalizationURL(rawURL))

	return shortURL, nil
}

// Ping PostgreSQL
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("pgx", h.config.DBDSN)
	if err != nil {
		http.Error(w, "DB connection failed", http.StatusInternalServerError)
		return
	}
	err = db.Ping()
	if err != nil {
		http.Error(w, "DB connection failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
