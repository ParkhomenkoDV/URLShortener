package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewURLManager(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		wantErr  bool
		checkURL bool
	}{
		{"Valid URL", "http://localhost:8080", false, true},
		{"Missing port", "http://localhost", true, false},
		{"Invalid URL", "not_a_url", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewURLManager(tt.baseURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewURLManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && manager.baseURL != tt.baseURL {
				t.Errorf("Expected baseURL %s, got %s", tt.baseURL, manager.baseURL)
			}
		})
	}
}

func TestShortenHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
	}{
		{"Valid POST", http.MethodPost, "https://google.com", http.StatusCreated},
		{"Invalid method", http.MethodGet, "", http.StatusBadRequest},
		{"Empty body", http.MethodPost, "", http.StatusBadRequest},
		{"Invalid URL", http.MethodPost, "ftp://test", http.StatusBadRequest},
	}
	manager, _ := NewURLManager("http://localhost:8080")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			manager.Shorten(w, req)

			res := w.Result()
			if res.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, res.StatusCode)
			}

			if tt.wantStatus == http.StatusCreated {
				if !strings.HasPrefix(w.Body.String(), "http://localhost:8080/") {
					t.Errorf("Expected shortened URL, got %s", w.Body.String())
				}
			}
		})
	}
}

func TestExpandHandler(t *testing.T) {
	manager, _ := NewURLManager("http://localhost:8080")
	manager.urls["test123"] = "https://google.com"

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantLoc    string
	}{
		{"Valid GET", http.MethodGet, "/test123", http.StatusTemporaryRedirect, "https://google.com"},
		{"Invalid method", http.MethodPost, "/test123", http.StatusBadRequest, ""},
		{"Not found", http.MethodGet, "/missing", http.StatusBadRequest, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			manager.Expand(w, req)

			res := w.Result()
			if res.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, res.StatusCode)
			}

			if tt.wantLoc != "" {
				loc := res.Header.Get("Location")
				if loc != tt.wantLoc {
					t.Errorf("Expected Location %s, got %s", tt.wantLoc, loc)
				}
			}
		})
	}
}

/*
func TestIntegration(t *testing.T) {
	manager, _ := NewURLManager("http://localhost:8080")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			manager.Shorten(w, r)
		case http.MethodGet:
			manager.Expand(w, r)
		}
	}))
	defer ts.Close()

	// Test shortening
	origURL := "https://example.com"
	resp, err := http.Post(ts.URL, "text/plain", strings.NewReader(origURL))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	shortURL, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Test expanding
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err = client.Get(string(shortURL))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("Expected status 307, got %d", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	if loc != origURL {
		t.Errorf("Expected Location %s, got %s", origURL, loc)
	}
}
*/

func Test_generateID(t *testing.T) {
	tests := []struct {
		name    string
		urls    map[string]string
		wantErr bool
	}{
		{"1", map[string]string{"https://test.com": ""}, false},
		{"2", map[string]string{"https://gmail.com": ""}, false},
		{"3", map[string]string{"https://yandex.ru": ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := generateID(tt.urls)
			if _, ok := tt.urls[id]; ok {
				t.Errorf("Generated non-unique ID %v", id)
			}
		})
	}
}
