package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type URLManager struct {
	mu      sync.Mutex
	urls    map[string]string // short: long
	baseURL string
	host    string
	port    string
}

func NewURLManager(baseURL string) (*URLManager, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return &URLManager{}, fmt.Errorf("incorrect baseURL")
	}

	if url.Port() == "" {
		return &URLManager{}, fmt.Errorf("incorrect baseURL port")
	}

	return &URLManager{
		urls:    make(map[string]string),
		baseURL: baseURL,
		host:    url.Host,
		port:    url.Port(),
	}, nil
}

// Shorten - преобразователь длинного url в короткий.
func (manager *URLManager) Shorten(w http.ResponseWriter, r *http.Request) {
	originalURL, err := readBody(r)
	if err != nil || len(originalURL) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	manager.mu.Lock()
	id := generateID(manager.urls)
	manager.mu.Unlock()

	manager.urls[id] = originalURL

	shortenedURL := fmt.Sprintf("http://localhost:8080/%s", id)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(shortenedURL))
}

// Expand - преобразователь короткого url в длинный.
func (manager *URLManager) Expand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")

	manager.mu.Lock()
	originalURL, exists := manager.urls[id]
	manager.mu.Unlock()

	if !exists {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func readBody(r *http.Request) (string, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func generateID(urls map[string]string) string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		log.Fatal(err)
	}
	id := base64.URLEncoding.EncodeToString(b)

	//id = strings.ReplaceAll(id, "=", "c")
	//id = strings.ReplaceAll(id, "_", "D")
	//id = strings.ReplaceAll(id, "-", "G")

	if _, exists := urls[id]; exists {
		return generateID(urls)
	}

	return id
}

func main() {
	manager, err := NewURLManager("http://localhost:8080")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			manager.Shorten(w, r)
		case http.MethodGet:
			manager.Expand(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusBadRequest)
		}
	})

	log.Printf("Server starting on %s \n", manager.baseURL)
	log.Fatal(http.ListenAndServe(":"+manager.port, nil))
}
