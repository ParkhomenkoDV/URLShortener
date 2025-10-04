package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/storage"
)

type Server struct {
	config     *config.Config
	httpServer *http.Server
	db         *storage.DB
}

// New - создание нового сервера.
func New(config *config.Config, handler http.Handler, db *storage.DB) *Server {
	return &Server{
		config: config,
		httpServer: &http.Server{
			Addr:    config.ServerAddress,
			Handler: handler,
		},
		db: db,
	}
}

// Start - запуск сервера.
func (s *Server) Start() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server started at %s \n", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	<-quit // Ожидаем сигнал завершения
	log.Println("Server ended")

	return s.shutdown()
}

// shutdown - graceful shutdown.
func (s *Server) shutdown() error {
	// Создание контекста с таймаутом для shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Завершение работы HTTP сервера
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("error starting server: %v", err)
		return err
	}

	// Сохраняем данные перед завершением
	return s.saveData()
}

// saveData - сохраняет данные в файл.
func (s *Server) saveData() error {
	log.Println("Saving data")
	if err := s.db.SaveToFile(s.config.FileStorage); err != nil {
		log.Printf("saving data error: %v", err)
		return err
	}

	log.Println("Data saved successfully")
	return nil
}
