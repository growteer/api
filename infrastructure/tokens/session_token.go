package tokens

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (p *Provider) NewSessionToken(address string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject: address,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.sessionTTL)),
	})

	tokenString, err := token.SignedString(p.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create session token")
	}

	return tokenString, nil
}