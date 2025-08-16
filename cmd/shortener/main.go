package main

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
