package server

import (
	"flag"
	"fmt"
	"net/url"
)

const defaultAddress = "localhost:8080"

type Config struct {
	ServerAddress string
	BaseURL       string
}

func New() Config {
	config := Config{}

	flag.StringVar(&config.ServerAddress, "a", defaultAddress, "Server address")
	flag.StringVar(&config.BaseURL, "b", "", "Base URL")
	flag.Parse()

	if config.ServerAddress == "" {
		config.ServerAddress = defaultAddress
	}

	if config.BaseURL == "" {
		config.BaseURL = "http://" + defaultAddress
	}

	if _, err := url.ParseRequestURI(config.BaseURL); err != nil {
		panic(fmt.Sprintf("Invalid base URL: %s", config.BaseURL))
	}

	return config
}
