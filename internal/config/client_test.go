package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientConfig_Default(t *testing.T) {
	resetFlags()

	os.Unsetenv("HOST")

	cfg := NewClientConfig()

	assert.Equal(t, DefaultHost+DefaultPort, cfg.Host)
}

func TestNewClientConfig_WithEnvVar(t *testing.T) {
	resetFlags()

	os.Setenv("HOST", "http://192.168.1.1:7070")
	defer os.Unsetenv("HOST")

	cfg := NewClientConfig()

	assert.Equal(t, "http://192.168.1.1:7070", cfg.Host)
}
