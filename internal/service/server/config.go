package server

import (
	"flag"
	"net/url"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func New() (Config, error) {
	config := Config{}

	flag.StringVar(&config.ServerAddress, "a", config.ServerAddress, "Server address")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "Base URL")
	flag.Parse()

	err := env.Parse(&config)
	if err != nil {
		return config, err
	}

	if _, err := url.ParseRequestURI(config.BaseURL); err != nil {
		return config, err
	}

	return config, nil
}
