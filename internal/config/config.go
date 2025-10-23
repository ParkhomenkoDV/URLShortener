package config

// Конфигурация проекта
type Config struct {
	Protocol      string
	Port          string
	ShortAddress  string
	FilePath      string `doc:"Путь к локальной директории БД"`
	AddressDB     string `doc:"Адрес БД"`
	AuthSecretKey string
}

// Инициализация конфигурации
func New() *Config {
	// Парсинг флагов
	reqAddress, resAddress, filePath, adressDB := parseFlags()

	return &Config{
		Protocol:      "http://",
		Port:          reqAddress,
		ShortAddress:  resAddress,
		FilePath:      filePath,
		AddressDB:     adressDB,
		AuthSecretKey: "your-secret-key-change-in-production", // В production брать из переменной окружения
	}
}
