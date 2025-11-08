package config

// Config представляет основную структуру конфигурации приложения.
// Содержит все необходимые параметры для работы сервера сокращения URL.
type Config struct {
	Protocol      string // Протокол HTTP/HTTPS
	Port          string // Порт сервера в формате ":8080"
	ShortAddress  string // Базовый адрес для сокращенных URL
	FilePath      string `doc:"Путь к локальной директории БД"` // Путь к файлу для хранения данных
	AddressDB     string `doc:"Адрес БД"`                       // DSN строка для подключения к базе данных
	AuthSecretKey string // Секретный ключ для аутентификации
	LengthKey     int    `doc:"Длина короткого ID"` // Длина генерируемых коротких идентификаторов
}

// New создает и инициализирует новую конфигурацию приложения.
// Приоритет параметров:
//  1. Переменные окружения
//  2. Флаги командной строки
//  3. Значения по умолчанию
//
// Возвращает готовую к использованию конфигурацию.
func New() *Config {
	cfg := parseFlags() // Парсинг флагов командной строки и переменных окружения

	return &Config{
		Protocol:      "http://",
		Port:          cfg.Port,
		ShortAddress:  cfg.ResAddress,
		FilePath:      cfg.FilePath,
		AddressDB:     cfg.AddressDB,
		AuthSecretKey: "your-secret-key-change-in-production", // TODO: В production брать из переменной окружения
		LengthKey:     6,                                      // Стандартная длина для коротких URL
	}
}
