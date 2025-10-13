package config

// Конфигурация проекта
type Config struct {
	Protocol     string
	Port         string
	ShortAddress string
	FilePath     string
	AddressDB    string
}

// Инициализация конфигурации
func New() *Config {
	// Получение данных из флагов
	reqAddress, resAddress, filePath, adressDB := parseFlags()

	return &Config{
		Protocol:     "http://",
		Port:         reqAddress,
		ShortAddress: resAddress,
		FilePath:     filePath,
		AddressDB:    adressDB,
	}
}
