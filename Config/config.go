package Config

import (
	"os"
)

type ServerConfig struct {
	IP   string
	Port string
}

func GetServerConfig() ServerConfig {
	ip := os.Getenv("SERVER_IP")
	if ip == "" {
		ip = "0.0.0.0" // значение по умолчанию
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // значение по умолчанию
	}

	return ServerConfig{
		IP:   ip,
		Port: port,
	}
}
