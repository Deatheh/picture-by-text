package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	jwt.MapClaims
	UUID string `json:"user_id"`
	Type string `json:"type"`
}

func GenerateAccessToken(UUID string, expirationTime int, signingKey string) (string, error) {
	if UUID == "" {
		return "", errors.New("UUID is empty")
	}
	claims := &TokenClaims{
		jwt.MapClaims{
			"ExpiresAt": time.Now().Add(time.Duration(expirationTime) * time.Minute).Unix(),
			"IssuedAr":  time.Now().Unix(),
		},
		UUID,
		"access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingKey))
}

func GenerateRefreshToken(UUID string, expirationTime int, signingKey string) (string, error) {
	if UUID == "" {
		return "", errors.New("UUID is empty")
	}
	claims := &TokenClaims{
		jwt.MapClaims{
			"ExpiresAt": time.Now().Add(time.Duration(expirationTime) * time.Hour).Unix(),
			"IssuedAr":  time.Now().Unix(),
		},
		UUID,
		"refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingKey))
}

func ParseToken(tokenString, signingKey string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, err
	}

	if time.Now().Unix() > int64(claims.MapClaims["ExpiresAt"].(float64)) {
		return nil, errors.New("Token has expired")
	}

	return claims, nil
}
