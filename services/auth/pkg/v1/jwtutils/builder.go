package jwtutils

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type Builder struct {
	secret string
	exp    time.Duration
	method jwt.SigningMethod
	claims map[string]interface{}
}

func NewBuilder(secret string, exp time.Duration, method jwt.SigningMethod) *Builder {
	return &Builder{
		secret: secret,
		exp:    exp,
		method: method,
		claims: make(map[string]interface{}),
	}
}

func (b *Builder) Claim(key string, value interface{}) *Builder {
	b.claims[key] = value

	return b
}

func (b *Builder) BuildRaw() (string, error) {
	b.claims["exp"] = time.Now().Add(b.exp).Unix()

	token := jwt.NewWithClaims(b.method, jwt.MapClaims(b.claims))

	return token.SignedString([]byte(b.secret))
}
