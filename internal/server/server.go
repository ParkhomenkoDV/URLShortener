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

// Server представляет HTTP сервер приложения с поддержкой graceful shutdown
type Server struct {
	httpServer *http.Server     // Встроенный HTTP сервер
	service    *service.Service // Сервисный слой приложения
}

// New создает новый экземпляр сервера с указанными параметрами
func New(address string, handler http.Handler, service *service.Service) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    address, // Адрес для прослушивания (например: ":8080")
			Handler: handler, // HTTP обработчик (роутер)
		},
		service: service, // Сервис для бизнес-логики
	}
}

// Start запускает HTTP сервер и обрабатывает graceful shutdown
// Метод блокирует выполнение до получения сигнала завершения
func (s *Server) Start() error {
	// Канал для получения сигналов ОС о завершении работы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		log.Printf("Server started at %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup error: %v", err)
		}
	}()

	// Ожидаем сигнал завершения от ОС
	<-quit
	log.Println("Shutting down the server")

	// Выполняем graceful shutdown
	return s.shutdown()
}

// shutdown выполняет корректное завершение работы сервера
// Закрывает HTTP соединения и освобождает ресурсы
func (s *Server) shutdown() error {
	// Создаём контекст с таймаутом для graceful shutdown
	// Сервер имеет 30 секунд чтобы завершить активные соединения
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Останавливаем HTTP сервер, позволяя завершить активные запросы
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server termination error: %v", err)
		return err
	}

	// Закрываем соединения с хранилищем данных
	return s.closeRepository()
}

// closeRepository освобождает ресурсы, связанные с хранилищем данных
func (s *Server) closeRepository() error {
	log.Println("Closing the connection to the repository")
	if err := s.service.Close(); err != nil {
		log.Printf("Closing repository error: %v", err)
		return err
	}

	log.Println("Repository connection closed successfully")
	return nil
}
