package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenManager struct {
	jwtSignature string
}

func NewTokenManager(jwtSignature string) TokenManager {
	return TokenManager{
		jwtSignature: jwtSignature,
	}
}

func (tm *TokenManager) GenerateJWT() (string, error) {
	signKey := []byte(tm.jwtSignature)
	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["exp"] = now.Add(1 * time.Hour).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signKey)

	if err != nil {
		return "", fmt.Errorf("failed to create a signed token %w", err)
	}

	return token, nil
}

func (tm *TokenManager) ValidateToken(tokenString string) (bool, error) {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// on success return our secret to satisfy the parse function
		// if the signature on token != our returned signature, returns error
		return []byte(tm.jwtSignature), nil
	})
	if err != nil {
		return false, fmt.Errorf("failed to validate the token: %v", err)
	}

	return true, nil
}
