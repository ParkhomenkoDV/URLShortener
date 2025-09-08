package main

import (
	"log"
	"net/http"

	"github.com/ParkhomenkoDV/URLShortener/internal/handler"
	"github.com/ParkhomenkoDV/URLShortener/internal/logger"
	"github.com/ParkhomenkoDV/URLShortener/internal/middleware"
	"github.com/ParkhomenkoDV/URLShortener/internal/service/server"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	config, err := server.New()
	if err != nil {
		panic(err)
	}

	logger.New()

	handler := handler.New(config)

	r := chi.NewRouter()
	r.Use(middleware.GzipRequestMiddleware)
	r.Use(middleware.GzipResponseMiddleware)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(logger.LoggingMiddleware)

	r.Post("/", handler.Post)
	r.Post("/api/shorten", handler.PostJSON)
	r.Get("/{id}", handler.Get)

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	log.Printf("Server started at %s \n", config.ServerAddress)
	log.Fatal(http.ListenAndServe(config.ServerAddress, r))
}
