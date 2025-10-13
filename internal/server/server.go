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

// Представляет HTTP сервер с graceful shutdown
type Server struct {
	httpServer *http.Server
	service    *service.Service
}

// Создаёт новый сервер
func New(addr string, handler http.Handler, service *service.Service) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		service: service,
	}
}

// Запускает сервер с graceful shutdown
func (s *Server) Start() error {
	// Канал для получения сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Сервер запущен на %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	<-quit
	log.Println("Завершение работы сервера...")

	return s.shutdown()
}

// graceful shutdown
func (s *Server) shutdown() error {
	// Создаём контекст с таймаутом для shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Завершаем HTTP сервер
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при завершении сервера: %v", err)
		return err
	}

	// Закрываем соединение с репозиторием
	return s.closeRepository()
}

// Закрывает соединение с репозиторием
func (s *Server) closeRepository() error {
	log.Println("Закрытие соединения с репозиторием...")
	if err := s.service.Close(); err != nil {
		log.Printf("Ошибка при закрытии репозитория: %v", err)
		return err
	}

	log.Println("Соединение с репозиторием успешно закрыто")
	return nil
}
