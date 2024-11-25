package config

import (
	"flag"
	"os"
)

// ServerConfig хранит информацию о конфгурации сервера.
type ServerConfig struct {
	// Host - адрес эндпоинта HTTP-сервера.
	Host string `json:"host"`
	// Port
	Port string `json:"port"`
	// DatabaseDsn - строка с адресом подключения к БД.
	DatabaseDsn string `json:"database_dsn"`
	// JwtKey - ключ для хеша.
	JwtKey string
	// CryptoKey - путь до файла с публичным ключом.
	CryptoKey string `json:"crypto_key"`
	// CryptoCrt - путь до файла с сертификатом.
	CryptoCrt string `json:"crypto_crt"`
}

// NewServerConfig создаёт и возвращает новый экземпляр ServerConfig, инициализированный с помощью флагов.
func NewServerConfig() ServerConfig {
	sc := ServerConfig{}
	sc.parseFlag()
	return sc
}

func (sc *ServerConfig) parseFlag() {
	var scFlags ServerConfig
	flag.StringVar(&scFlags.Host, "host", DefaultHost, "host to run server")
	flag.StringVar(&scFlags.Port, "port", DefaultPort, "port to run server")
	flag.StringVar(&scFlags.DatabaseDsn, "dsn", DefaultDatabaseDSN, "database dsn")
	flag.StringVar(&scFlags.JwtKey, "jwt-key", DefaultJwtKey, "jwt key")
	flag.StringVar(&scFlags.CryptoKey, "crypto-key", DefaultCryptoKeyServer, "crypto key")
	flag.StringVar(&scFlags.CryptoCrt, "crypto-crt", DefaultCryptoCrtServer, "crypto crt")
	flag.Parse()

	if scFlags.Host != "" {
		sc.Host = scFlags.Host
	}
	if scFlags.Port != "" {
		sc.Port = scFlags.Port
	}
	if scFlags.DatabaseDsn != "" {
		sc.DatabaseDsn = scFlags.DatabaseDsn
	}
	if scFlags.JwtKey != "" {
		sc.JwtKey = scFlags.JwtKey
	}
	if scFlags.CryptoKey != "" {
		sc.CryptoKey = scFlags.CryptoKey
	}
	if scFlags.CryptoCrt != "" {
		sc.CryptoCrt = scFlags.CryptoCrt
	}

	if envHost := os.Getenv("HOST"); envHost != "" {
		sc.Host = envHost
	}
	if envPort := os.Getenv("PORT"); envPort != "" {
		sc.Port = envPort
	}
	if envDatabaseDsn := os.Getenv("DATABASE_DSN"); envDatabaseDsn != "" {
		sc.DatabaseDsn = envDatabaseDsn
	}
	if envJwtKey := os.Getenv("JWT_KEY"); envJwtKey != "" {
		sc.JwtKey = envJwtKey
	}
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		sc.CryptoKey = envCryptoKey
	}
	if envCryptoCrt := os.Getenv("CRYPTO_CRT"); envCryptoCrt != "" {
		sc.CryptoCrt = envCryptoCrt
	}
}
