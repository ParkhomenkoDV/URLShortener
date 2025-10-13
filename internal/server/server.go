package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ParkhomenkoDV/URLShortener/internal/service"
)

// HTTP сервер с graceful shutdown
type Server struct {
	httpServer *http.Server
	service    *service.Service
}

// Конструктор сервера
func New(address string, handler http.Handler, service *service.Service) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: handler,
		},
		service: service,
	}
}

// Запускает сервер
func (s *Server) Start() error {
	// Канал для получения сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Server started at %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup error: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	<-quit

	log.Println("Shutting down the server")

	return s.shutdown()
}

// graceful shutdown
func (s *Server) shutdown() error {
	// Создаём контекст с таймаутом для shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Завершаем HTTP сервер
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server termination error: %v", err)
		return err
	}

	// Закрываем соединение с репозиторием
	return s.closeRepository()
}

// Закрывает соединение с репозиторием
func (s *Server) closeRepository() error {
	log.Println("Closing the connection to the repository")
	if err := s.service.Close(); err != nil {
		log.Printf("Closing repository error: %v", err)
		return err
	}

	log.Println("Repository connection closed.")
	return nil
}
