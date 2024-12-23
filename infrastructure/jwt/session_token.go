package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenProvider struct {
	secretKey []byte
	sessionTTL time.Duration
}

func NewTokenProvider(secretKey string, sessionTTLMinutes int) *TokenProvider {
	return &TokenProvider{
		secretKey: []byte(secretKey),
		sessionTTL: time.Minute * time.Duration(sessionTTLMinutes),
	}
}

func (p *TokenProvider) NewSessionToken(address string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": address,
		"exp": time.Now().Add(p.sessionTTL).Unix(),
	})

	tokenString, err := token.SignedString(p.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create token")
	}

	return tokenString, nil
}