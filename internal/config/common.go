package config

// Конфигурация по умолчанию.
const (
	DefaultHost            = "localhost"
	DefaultJwtKey          = "popa"
	DefaultDatabaseDSN     = "host=localhost port=5432 user=mila dbname=goph_keeper password=qwerty sslmode=disable"
	DefaultPort            = ":8080"
	DefaultCryptoKeyServer = "./crypt/server.key"
	DefaultCryptoCrtServer = "./crypt/server.crt"
)
