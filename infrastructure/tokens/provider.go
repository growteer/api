package tokens

import (
	"time"
)

type Provider struct {
	secretKey []byte
	sessionTTL time.Duration
	refreshTTL time.Duration
}

func NewProvider(secretKey string, sessionTTLMinutes, refreshTTLMinutes int) *Provider {
	return &Provider{
		secretKey: []byte(secretKey),
		sessionTTL: time.Minute * time.Duration(sessionTTLMinutes),
		refreshTTL: time.Minute * time.Duration(refreshTTLMinutes),
	}
}
