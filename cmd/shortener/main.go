package main

import (
	"log"
	"net/http"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/handler"
	"github.com/ParkhomenkoDV/URLShortener/internal/logger"
	"github.com/ParkhomenkoDV/URLShortener/internal/middleware"
	"github.com/ParkhomenkoDV/URLShortener/internal/server"
	"github.com/ParkhomenkoDV/URLShortener/internal/storage"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	logger.New()

	db := storage.New()
	db.LoadFromFile(cfg.FileStorage)

	hand := handler.New(cfg, db)

	r := chi.NewRouter()
	r.Use(middleware.GzipRequestMiddleware)
	r.Use(middleware.GzipResponseMiddleware)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(logger.LoggingMiddleware)

	r.Post("/", hand.Post)
	r.Post("/api/shorten", hand.PostJSON)
	r.Get("/{id}", hand.Get)

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	srv := server.New(&cfg, r, db)
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}

	log.Println("Server ended")
}
