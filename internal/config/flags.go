package config

import (
	"flag"
	"log"
	"os"
	"strings"
)

// parseFlags обрабатывает аргументы командной строки
func parseFlags() (string, string, string, string) {
	// адрес запуска HTTP-сервера значением localhost:8080 по умолчанию
	portFlag := flag.String("a", "localhost:8080", "address and port to run server")

	// базовый адрес результирующего сокращённого URL значением
	resAddressFlag := flag.String("b", "http://localhost:8080", "address and port for short url")

	// путь к файлу с урлами, значение data/urls.json по умолчанию
	filePathFlag := flag.String("f", "data/urls.json", "path to the file for storing data")

	// адрес для базы данных
	addressFlagDB := flag.String("d", "", "database address")

	flag.Parse()

	// Проверка и обработка адреса сервера
	portParts := strings.Split(*portFlag, ":")
	if len(portParts) < 2 {
		log.Printf("invalid address format: %s, expected format: host:port\n", *portFlag)
		flag.Usage()
		os.Exit(2)
	}

	port := ":" + portParts[1]

	// Проверка базового адреса
	resAddress := *resAddressFlag
	if !strings.HasPrefix(resAddress, "http://") && !strings.HasPrefix(resAddress, "https://") {
		log.Printf("invalid base address: %s, must start with http:// or https://\n", resAddress)
		flag.Usage()
		os.Exit(2)
	}

	// Путь до файла с urls
	filePath := *filePathFlag

	// Адрес для базы данных
	addressDB := *addressFlagDB

	// Приоритет параметров согласно заданию:
	// 1. Переменная окружения (наивысший приоритет)
	// 2. Флаг командной строки
	// 3. Значение по умолчанию (уже установлено)

	// Если параметры заданы через переменные окружения, используем их
	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		// Обрабатываем переменную окружения SERVER_ADDRESS
		envPortParts := strings.Split(envServerAddr, ":")
		if len(envPortParts) >= 2 {
			port = ":" + envPortParts[1]
		} else {
			port = envServerAddr
		}
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		resAddress = envBaseURL
	}

	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		filePath = envFilePath
	}

	if envAddressDB := os.Getenv("DATABASE_DSN"); envAddressDB != "" {
		addressDB = envAddressDB
	}

	return port, resAddress, filePath, addressDB
}
