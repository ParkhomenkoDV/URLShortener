package main

import (
	"log"
	"net/http"

	"github.com/ParkhomenkoDV/URLShortener/internal/handler"
	"github.com/ParkhomenkoDV/URLShortener/internal/service/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config := server.New()

	handler := handler.New(config)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", handler.Post)
	r.Post("/api/shorten", handler.PostJSON) // Новый endpoint для JSON
	r.Get("/{id}", handler.Get)

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	log.Printf("Server started at %s \n", config.ServerAddress)
	log.Fatal(http.ListenAndServe(config.ServerAddress, r))
}
