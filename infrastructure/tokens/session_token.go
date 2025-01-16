package tokens

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/growteer/api/pkg/web3util"
)

func (p *Provider) NewSessionToken(did *web3util.DID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject: did.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.sessionTTL)),
	})

	tokenString, err := token.SignedString(p.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create session token")
	}

	return tokenString, nil
}