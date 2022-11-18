package util

import (
	"fmt"
	"time"

	authpb "github.com/MyyPo/grpc-chat/pb/auth/v1"
	"github.com/golang-jwt/jwt/v4"
)

type TokenManager struct {
	AccessSignature      string
	RefreshSignature     string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func NewTokenManager(accessSignature, refreshSignature string, accessTokenDuration, refreshTokenDuration time.Duration) TokenManager {
	return TokenManager{
		AccessSignature:      accessSignature,
		RefreshSignature:     refreshSignature,
		AccessTokenDuration:  accessTokenDuration,
		RefreshTokenDuration: refreshTokenDuration,
	}
}

func (tm *TokenManager) GenerateJWT(isAccessTok bool) (string, error) {
	var signature string
	var duration time.Duration

	// add different claims based on the token type
	if isAccessTok {
		signature = tm.AccessSignature
		duration = tm.AccessTokenDuration
	} else {
		signature = tm.RefreshSignature
		duration = tm.RefreshTokenDuration
	}

	signKey := []byte(signature)
	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["exp"] = now.Add(duration).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signKey)

	if err != nil {
		return "", fmt.Errorf("failed to create a signed token %w", err)
	}

	return token, nil
}

func (tm *TokenManager) ValidateToken(tokenString string, isAccessTok bool) error {
	var signature string
	if isAccessTok {
		signature = tm.AccessSignature
	} else {
		signature = tm.RefreshSignature
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// on success return our secret to satisfy the parse function
		// if the signature on token != our returned signature, returns error
		return []byte(signature), nil
	})
	if err != nil {
		return fmt.Errorf("failed to validate the token: %v", err)
	}

	// check if the token expired, or something else is wrong with claims
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("claims error: %v", err)
	}

	return nil
}

func (tm *TokenManager) RerfreshToken(tokenString string) (authpb.RefreshTokenResponse, error) {
	if err := tm.ValidateToken(tokenString, false); err != nil {
		return authpb.RefreshTokenResponse{}, err
	}

	accessToken, _ := tm.GenerateJWT(true)
	refreshToken, _ := tm.GenerateJWT(false)

	return authpb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
