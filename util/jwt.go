package util

import "github.com/golang-jwt/jwt/v4"

func GenerateJWT(key string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	signKey := []byte(key)
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
