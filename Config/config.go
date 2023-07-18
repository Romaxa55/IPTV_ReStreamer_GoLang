package Config

import (
	"os"
)

type ServerConfig struct {
	IP      string
	Port    string
	IptvUrl string
	Token   string
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

	iptv := os.Getenv("IPTV_URL")
	if iptv == "" {
		iptv = "http://4dc7c7fa1a76.faststreem.org/playlists/uplist/42270cf01929b705c42612e12df9ef9f/playlist.m3u8" // значение по умолчанию
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		token = "DefaultPassword" // значение по умолчанию
	}

	return ServerConfig{
		IP:      ip,
		Port:    port,
		IptvUrl: iptv,
		Token:   token,
	}
}
