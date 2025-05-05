package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenProvider struct {
	secret string
}

func NewTokenProvider(secret string) *TokenProvider {
	return &TokenProvider{secret: secret}
}

func (p *TokenProvider) Generate(email string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	return t.SignedString([]byte(p.secret))
}

func (p *TokenProvider) Validate(tokenStr string) bool {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.secret), nil
	})
	return err == nil && token.Valid
}
