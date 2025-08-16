package main

import (
	"log"
	"net/http"

	"github.com/ParkhomenkoDV/URLShortener/internal/urlmanager"
	"github.com/go-chi/chi/v5"
)

func main() {
	manager, err := urlmanager.New("http://localhost:8080")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	r.Post("/", manager.Shorten)

	r.Get("/{id}", manager.Expand)

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	log.Printf("Server starting on %s \n", manager.URL)
	log.Fatal(http.ListenAndServe(":"+manager.Port, r))
}

/*
import (
	"log"
	"net/http"

	"github.com/ParkhomenkoDV/URLShortener/internal/urlmanager"
)

func main() {
	manager, err := urlmanager.New("http://localhost:8080")
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

	log.Printf("Server starting on %s \n", manager.URL)
	log.Fatal(http.ListenAndServe(":"+manager.Port, nil))
}
*/
