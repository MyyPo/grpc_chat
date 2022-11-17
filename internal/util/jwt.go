package util

import "github.com/golang-jwt/jwt/v4"

type TokenManager struct {
	jwtSignature string
}

func NewTokenManager(jwtSignature string) TokenManager {
	return TokenManager{
		jwtSignature: jwtSignature,
	}
}

func (t *TokenManager) GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	signKey := []byte(t.jwtSignature)
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
