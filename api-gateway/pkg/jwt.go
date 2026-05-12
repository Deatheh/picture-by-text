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
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Проверка срока действия
	if exp, ok := claims.MapClaims["ExpiresAt"]; ok {
		expTime := time.Unix(int64(exp.(float64)), 0)
		if time.Now().After(expTime) {
			return nil, errors.New("token expired")
		}
	}

	return claims, nil
}
