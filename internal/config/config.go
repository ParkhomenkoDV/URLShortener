package config

import (
	"flag"
	"net/url"

	"github.com/caarlos0/env/v11"
)

const defaultAddress = "localhost:8080"

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"` // envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL"`       // envDefault:"http://localhost:8080"`
	FileStorage   string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() (Config, error) {
	config := Config{}

	err := env.Parse(&config)
	if err != nil {
		return config, err
	}

	configFlags := Config{}

	flag.StringVar(&configFlags.ServerAddress, "a", defaultAddress, "Server address")
	flag.StringVar(&configFlags.BaseURL, "b", "http://"+defaultAddress, "Base URL")
	flag.StringVar(&configFlags.FileStorage, "f", "data/db.json", "File Storage")
	flag.Parse()

	if config.ServerAddress == "" {
		config.ServerAddress = configFlags.ServerAddress
	}
	if config.BaseURL == "" {
		config.BaseURL = configFlags.BaseURL
	}
	if config.FileStorage == "" {
		config.FileStorage = configFlags.FileStorage
	}

	if _, err := url.ParseRequestURI(config.BaseURL); err != nil {
		return config, err
	}

	return config, nil
}
