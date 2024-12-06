package github

import (
	"crypto/rsa"
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	// gh "github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
)

const (
	// AppTokenExpiration is the default expiration time for the GitHub App token.
	AppTokenExpiration = 10 * time.Minute // 10 minutes (maximum allowed by GitHub)
	bearerTokenType    = "Bearer"
)

type appTokenSource struct {
	clientID   string
	privateKey *rsa.PrivateKey
	expiration time.Duration
}

// NewApplicationTokenSource creates a new GitHub App token source using the provided
// client ID and private key.
func NewApplicationTokenSource(clientID string, privateKey []byte) (oauth2.TokenSource, error) {
	if clientID == "" {
		return nil, errors.New("clientID is required")
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}

	t := &appTokenSource{
		clientID:   clientID,
		privateKey: privKey,
		expiration: AppTokenExpiration,
	}

	return t, nil
}

// Token generates a new GitHub App (JWT) token for authenticating as a GitHub App.
func (t *appTokenSource) Token() (*oauth2.Token, error) {
	now := time.Now().Add(-60 * time.Second) // 1 minute in the past
	expiresAt := now.Add(t.expiration)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Issuer:    t.clientID,
	})

	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken: tokenString,
		TokenType:   bearerTokenType,
		Expiry:      expiresAt,
	}, nil
}
