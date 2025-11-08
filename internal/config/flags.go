package config

import (
	"flag"
	"log"
	"os"
	"strings"
)

// ConfigFlags структура для временного хранения конфигурационных параметров
// перед их преобразованием в основную структуру Config
type ConfigFlags struct {
	Port       string // Порт сервера
	ResAddress string // Базовый адрес для сокращенных URL
	FilePath   string // Путь к файлу хранилища
	AddressDB  string // Адрес базы данных
}

// parseFlags обрабатывает аргументы командной строки и переменные окружения.
// Приоритет получения параметров (от высшего к низшему):
//  1. Переменные окружения
//  2. Флаги командной строки
//  3. Значения по умолчанию
//
// Возвращает ConfigFlags с обработанными значениями.
func parseFlags() ConfigFlags {
	var cfg ConfigFlags

	// Создаем новый FlagSet для изоляции парсинга и избежания конфликтов в тестах
	fs := flag.NewFlagSet("config", flag.ContinueOnError)

	// Определение флагов командной строки с значениями по умолчанию
	portFlag := fs.String("a", "localhost:8080", "HTTP server address and port (format: host:port)")
	resAddressFlag := fs.String("b", "http://localhost:8080", "Base address for shortened URLs (must include http:// or https://)")
	filePathFlag := fs.String("f", "data/urls.json", "Path to local file storage for URLs")
	addressDBFlag := fs.String("d", "", "Database connection string (DSN)")

	// Парсим аргументы командной строки, если они присутствуют
	// В тестовой среде os.Args может быть пустым или содержать аргументы тестов
	if len(os.Args) > 1 {
		// Игнорируем ошибки парсинга, так как в тестах могут быть другие флаги
		_ = fs.Parse(os.Args[1:])
	}

	// Валидация и обработка адреса сервера
	portParts := strings.Split(*portFlag, ":")
	if len(portParts) < 2 {
		log.Printf("Invalid address format: %s, expected format: host:port\n", *portFlag)
		fs.Usage()
		os.Exit(2)
	}
	cfg.Port = ":" + portParts[1] // Сохраняем только порт в формате ":8080"

	// Валидация базового адреса для сокращенных URL
	resAddress := *resAddressFlag
	if !strings.HasPrefix(resAddress, "http://") && !strings.HasPrefix(resAddress, "https://") {
		log.Printf("Invalid base address: %s, must start with http:// or https://\n", resAddress)
		fs.Usage()
		os.Exit(2)
	}
	cfg.ResAddress = resAddress

	// Сохраняем путь к файлу хранилища
	cfg.FilePath = *filePathFlag

	// Сохраняем адрес базы данных (может быть пустым, если используется файловое хранилище)
	cfg.AddressDB = *addressDBFlag

	// ПЕРЕМЕННЫЕ ОКРУЖЕНИЯ ИМЕЮТ ВЫСШИЙ ПРИОРИТЕТ ЧЕМ ФЛАГИ КОМАНДНОЙ СТРОКИ

	// SERVER_ADDRESS - адрес и порт сервера
	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		envPortParts := strings.Split(envServerAddr, ":")
		if len(envPortParts) >= 2 {
			cfg.Port = ":" + envPortParts[1] // Извлекаем порт из формата "host:port"
		} else {
			cfg.Port = envServerAddr // Используем как есть, если порт не указан
		}
	}

	// BASE_URL - базовый адрес для сокращенных URL
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.ResAddress = envBaseURL
	}

	// FILE_STORAGE_PATH - путь к файлу хранилища URL
	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		cfg.FilePath = envFilePath
	}

	// DATABASE_DSN - строка подключения к базе данных
	if envAddressDB := os.Getenv("DATABASE_DSN"); envAddressDB != "" {
		cfg.AddressDB = envAddressDB
	}

	return cfg
}
