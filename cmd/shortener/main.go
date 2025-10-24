package main

import (
	"log"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/handler"
	"github.com/ParkhomenkoDV/URLShortener/internal/middleware"
	"github.com/ParkhomenkoDV/URLShortener/internal/repository"
	"github.com/ParkhomenkoDV/URLShortener/internal/server"
	"github.com/ParkhomenkoDV/URLShortener/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация конфигурации со всеми необходимыми параметрами работы приложения
	configuration := config.New()

	// Инициализация логгера
	middleware.InitLogger()

	// Инициализация роутера
	ginEngine := gin.Default() // под капотом автоматически настраивает базовые middleware (recovery, logger, static files)

	// Создание репозитория в зависимости от конфигурации
	repo := repository.CreateRepository(configuration.AddressDB, configuration.FilePath)

	// Создание сервиса
	srvc := service.New(repo, configuration)

	// Создание обработчика
	handler.New(ginEngine, srvc, configuration)

	// Создание и запуск сервера с graceful shutdown
	srv := server.New(configuration.Port, ginEngine, srvc)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server terminated with error: %v", err)
	}

	log.Println("Server is offline")
}
