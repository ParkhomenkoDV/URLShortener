package config

// Конфигурация проекта
type Config struct {
	Protocol      string
	Port          string
	ShortAddress  string
	FilePath      string `doc:"Путь к локальной директории БД"`
	AddressDB     string `doc:"Адрес БД"`
	AuthSecretKey string
	LengthKey     int `doc:"Длина короткого ID"`
}

// Инициализация конфигурации
func New() *Config {
	cfg := parseFlags() // Парсинг флагов
	return &Config{
		Protocol:      "http://",
		Port:          cfg.Port,
		ShortAddress:  cfg.ResAddress,
		FilePath:      cfg.FilePath,
		AddressDB:     cfg.AddressDB,
		AuthSecretKey: "your-secret-key-change-in-production", // В production брать из переменной окружения
		LengthKey:     6,
	}
}
