package tokens

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/growteer/api/pkg/web3util"
)

func (p *Provider) NewRefreshToken(did *web3util.DID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject: did.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.refreshTTL)),
	})

	tokenString, err := token.SignedString(p.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create refresh token")
	}

	return tokenString, nil
}

func (p *Provider) ParseRefreshToken(token string) (claims *jwt.RegisteredClaims, err error) {
	claims = &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return p.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	return claims, nil
}
