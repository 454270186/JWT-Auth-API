package domain

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TOKEN_DURATION = time.Hour

type Login struct {
	Username string `db:"username"`
	CustomerID sql.NullString `db:"customer_id"`
	Accounts sql.NullString `db:"account_numbers"`
	Role string `db:"role"`
}

func (l Login) GenerateToken() (*string, error) {
	var claims jwt.MapClaims
	if l.Accounts.Valid && l.CustomerID.Valid {
		claims = l.claimsForUser()
	} else {
		claims = l.claimsForAdmin()
	}

	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTokenAsString, err := tokens.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		log.Println("Failed while signing token")
		return nil, err
	}
	
	return &signedTokenAsString, nil
}

func (l Login) claimsForUser() jwt.MapClaims {
	accounts := strings.Split(l.Accounts.String, ",")
	return jwt.MapClaims{
		"customer_id": l.CustomerID.String,
		"role": l.Role,
		"username": l.Username,
		"accounts": accounts,
		"exp": time.Now().Add(TOKEN_DURATION).Unix(),
	}
}

func (l Login) claimsForAdmin() jwt.MapClaims {
	return jwt.MapClaims{
		"role": l.Role,
		"username": l.Username,
		"exp": time.Now().Add(TOKEN_DURATION).Unix(),
	}
}