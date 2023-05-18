package dto

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

const HMAC_SAMPLE_SECRET = "hmacSampleSecret"

type RefreshRequest struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// 1. Token is invalid
// 2. Token is valid but hasn't expired
func (r RefreshRequest) IsAccessTokenValid() error {
	_, err := jwt.Parse(r.AccessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(HMAC_SAMPLE_SECRET), nil
	})

	if errors.Is(err, jwt.ErrTokenExpired) {
		return jwt.ErrTokenExpired
	}

	return errors.New("token hasn't expired")
}

func (r RefreshRequest) GetUserName() (string, error) {
	token, _ := jwt.Parse(r.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(HMAC_SAMPLE_SECRET), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["username"].(string), nil
	}

	return "", errors.New("error while get username from claims")
}