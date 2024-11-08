package main

import (
	"errors"
	"os"
)

type config struct {
	server *ServerConfig
}

type ServerConfig struct {
	Addr string
}

func readConfigFromEnv() (*config, error) {
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		return nil, errors.New("SERVER_ADDR environment variable not set")
	}
	return &config{
		server: &ServerConfig{
			Addr: addr,
		},
	}, nil
}
