package domain

import (
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type AuthToken struct {
	Token        string `json:"jwt_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthToken(claims jwt.MapClaims, reClaims jwt.MapClaims) (AuthToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTokenAsString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		log.Println("Failed while signing token")
		return AuthToken{}, errors.New("failed token")
	}
	
	reToken := jwt.NewWithClaims(jwt.SigningMethodHS256, reClaims)
	signedReTokenAsString, err := reToken.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		log.Println("Failed while signing refresh token")
		return AuthToken{}, errors.New("failed refresh token")
	}

	return AuthToken{
		Token: signedTokenAsString,
		RefreshToken: signedReTokenAsString,
	}, nil
}
