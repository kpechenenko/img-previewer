package main

type config struct {
	server *ServerConfig
}

// ServerConfig конфигурация веб сервера превьювера изображений.
type ServerConfig struct {
	Addr string
}

func getDefaultConfig() *config {
	return &config{
		server: &ServerConfig{
			Addr: ":8080",
		},
	}
}
