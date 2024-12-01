package config

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerConfig_Default(t *testing.T) {
	// Сброс флагов перед тестом
	resetFlags()

	// Удаляем переменные окружения
	os.Unsetenv("HOST")
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_DSN")
	os.Unsetenv("JWT_KEY")
	os.Unsetenv("CRYPTO_KEY")
	os.Unsetenv("CRYPTO_CRT")

	// Вызываем вашу функцию
	cfg := NewServerConfig()

	// Проверяем значения по умолчанию
	assert.Equal(t, DefaultHost, cfg.Host)
	assert.Equal(t, DefaultPort, cfg.Port)
	assert.Equal(t, DefaultDatabaseDSN, cfg.DatabaseDsn)
	assert.Equal(t, DefaultJwtKey, cfg.JwtKey)
	assert.Equal(t, DefaultCryptoKeyServer, cfg.CryptoKey)
	assert.Equal(t, DefaultCryptoCrtServer, cfg.CryptoCrt)
}

func TestNewServerConfig_WithEnvVars(t *testing.T) {
	// Сброс флагов перед тестом
	resetFlags()

	// Устанавливаем переменные окружения
	os.Setenv("HOST", "env-host")
	os.Setenv("PORT", ":8081")
	os.Setenv("DATABASE_DSN", "env-dsn")
	os.Setenv("JWT_KEY", "env-jwt-key")
	os.Setenv("CRYPTO_KEY", "env-crypto-key")
	os.Setenv("CRYPTO_CRT", "env-crypto-crt")
	defer func() {
		os.Unsetenv("HOST")
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_DSN")
		os.Unsetenv("JWT_KEY")
		os.Unsetenv("CRYPTO_KEY")
		os.Unsetenv("CRYPTO_CRT")
	}()

	// Вызываем вашу функцию
	cfg := NewServerConfig()

	// Проверяем, что значения из переменных окружения установлены
	assert.Equal(t, "env-host", cfg.Host)
	assert.Equal(t, ":8081", cfg.Port)
	assert.Equal(t, "env-dsn", cfg.DatabaseDsn)
	assert.Equal(t, "env-jwt-key", cfg.JwtKey)
	assert.Equal(t, "env-crypto-key", cfg.CryptoKey)
	assert.Equal(t, "env-crypto-crt", cfg.CryptoCrt)
}

func resetFlags() {
	os.Clearenv()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	resetForTesting(nil)
}

// ResetForTesting clears all flag state and sets the usage function as directed.
// After calling ResetForTesting, parse errors in flag handling will not
// exit the program.
func resetForTesting(usage func()) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)
	flag.CommandLine.Usage = func() {
		_, err := fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		if err != nil {
			log.Fatal(err)
		}
		flag.PrintDefaults()
	}
	flag.Usage = usage
}
