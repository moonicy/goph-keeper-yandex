package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenGenerator struct {
	jwtKey []byte
}

func NewTokenGenerator(jwtKey string) (*TokenGenerator, error) {
	if jwtKey == "" {
		return nil, errors.New("jwtKey is empty")
	}
	return &TokenGenerator{
		jwtKey: []byte(jwtKey),
	}, nil
}

func (tg *TokenGenerator) GenerateToken(userID uint64) (string, error) {
	claims := &jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tg.jwtKey)
}
