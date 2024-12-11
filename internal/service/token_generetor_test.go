package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTokenGenerator_Success(t *testing.T) {
	jwtKey := "secret_key"
	tg, err := NewTokenGenerator(jwtKey)

	assert.NoError(t, err)
	assert.NotNil(t, tg)
	assert.Equal(t, []byte(jwtKey), tg.jwtKey)
}

func TestNewTokenGenerator_Error(t *testing.T) {
	tg, err := NewTokenGenerator("")

	assert.Error(t, err)
	assert.Nil(t, tg)
	assert.Equal(t, "jwtKey is empty", err.Error())
}

func TestGenerateToken_Success(t *testing.T) {
	jwtKey := "secret_key"
	userID := uint64(12345)

	tg, err := NewTokenGenerator(jwtKey)
	assert.NoError(t, err)
	assert.NotNil(t, tg)

	tokenString, err := tg.GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtKey), nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, parsedToken)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(userID), claims["user_id"])
}

func TestGenerateToken_InvalidKey(t *testing.T) {
	jwtKey := ""

	tg, err := NewTokenGenerator(jwtKey)
	assert.Error(t, err)
	assert.Nil(t, tg)
}
