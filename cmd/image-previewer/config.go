package main

type config struct {
	server *ServerConfig
}

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
