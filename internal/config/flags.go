package config

import (
	"flag"
	"log"
	"os"
	"strings"
)

// ConfigFlags структура для хранения всех конфигурационных параметров
type ConfigFlags struct {
	Port       string
	ResAddress string
	FilePath   string
	AddressDB  string
}

// parseFlags обрабатывает аргументы командной строки и возвращает структуру с конфигурацией
func parseFlags() ConfigFlags {
	var cfg ConfigFlags
	// Создаем новый FlagSet для избежания конфликтов в тестах
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	// адрес работы HTTP-сервера (localhost:8080 по умолчанию)
	portFlag := fs.String("a", "localhost:8080", "address and port to run server")
	// базовый адрес результирующего сокращённого URL значением
	resAddressFlag := fs.String("b", "http://localhost:8080", "address and port for short url")
	// путь к локальному файлу БД, значение data/urls.json по умолчанию
	filePathFlag := fs.String("f", "data/urls.json", "path to the file for storing data")
	// адрес для БД
	addressDBFlag := fs.String("d", "", "database address")

	// В тестах os.Args может быть пустым или содержать аргументы теста
	if len(os.Args) > 1 {
		fs.Parse(os.Args[1:])
	}

	// Проверка и обработка адреса сервера
	portParts := strings.Split(*portFlag, ":")
	if len(portParts) < 2 {
		log.Printf("invalid address format: %s, expected format: host:port\n", *portFlag)
		flag.Usage()
		os.Exit(2)
	}
	cfg.Port = ":" + portParts[1]

	// Проверка базового адреса
	resAddress := *resAddressFlag
	if !strings.HasPrefix(resAddress, "http://") && !strings.HasPrefix(resAddress, "https://") {
		log.Printf("invalid base address: %s, must start with http:// or https://\n", resAddress)
		flag.Usage()
		os.Exit(2)
	}
	cfg.ResAddress = resAddress

	// Путь до файла с urls
	cfg.FilePath = *filePathFlag

	// Адрес для базы данных
	cfg.AddressDB = *addressDBFlag

	// Приоритет параметров:
	// 1. Переменная окружения
	// 2. Флаг командной строки
	// 3. Дефолтное значение

	// Если параметры заданы через переменные окружения, используем их
	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		// Обрабатываем переменную окружения SERVER_ADDRESS
		envPortParts := strings.Split(envServerAddr, ":")
		if len(envPortParts) >= 2 {
			cfg.Port = ":" + envPortParts[1]
		} else {
			cfg.Port = envServerAddr
		}
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.ResAddress = envBaseURL
	}
	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		cfg.FilePath = envFilePath
	}
	if envAddressDB := os.Getenv("DATABASE_DSN"); envAddressDB != "" {
		cfg.AddressDB = envAddressDB
	}

	return cfg
}
