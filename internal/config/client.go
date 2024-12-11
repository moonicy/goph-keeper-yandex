package config

import (
	"flag"
	"os"
)

// ClientConfig хранит информацию о конфгурации клиента.
type ClientConfig struct {
	// Host - адрес эндпоинта HTTP-сервера.
	Host string `json:"host"`
}

// NewClientConfig создаёт и возвращает новый экземпляр ClientConfig, инициализированный с помощью флагов.
func NewClientConfig() ClientConfig {
	sc := ClientConfig{}
	sc.parseFlag()
	return sc
}

func (sc *ClientConfig) parseFlag() {
	var scFlags ClientConfig
	flag.StringVar(&scFlags.Host, "host", DefaultHost+DefaultPort, "address and port for connect to server")
	flag.Parse()

	if scFlags.Host != "" {
		sc.Host = scFlags.Host
	}

	if envHost := os.Getenv("HOST"); envHost != "" {
		sc.Host = envHost
	}
}
