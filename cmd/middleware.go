package main

import (
	"errors"
	"jwtAuth/domain"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ValidateToken validates the exp and sign of jwt token
func ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwtTokenString := ctx.Query("token")
		if jwtTokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing token",
			})
		} else {
			token, err := jwt.Parse(jwtTokenString, func(t *jwt.Token) (interface{}, error) {
				secretKey := []byte(domain.HMAC_SAMPLE_SECRET)

				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}

				return secretKey, nil
			})

			if token.Valid {
				log.Println("token is valid")
				ctx.Next()
			} else if errors.Is(err, jwt.ErrTokenExpired) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"errors": "token has expired",
				})
			}
		}
	}
}