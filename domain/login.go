package domain

import (
	"database/sql"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TOKEN_DURATION = 5 * time.Second
const REFRESH_TOKEN_DURATION = 24 * time.Hour

type Login struct {
	Username string `db:"username"`
	CustomerID sql.NullString `db:"customer_id"`
	Accounts sql.NullString `db:"account_numbers"`
	Role string `db:"role"`
}

func (l Login) GenerateToken() (*AuthToken, error) {
	claims := l.claimsForAccessToken()
	reClaims := l.claimsForRefresh()

	userToken, err := NewAuthToken(claims, reClaims)
	if err != nil {
		return nil, err
	}

	return &userToken, nil
}

func (l Login) claimsForAccessToken() jwt.MapClaims {
	if l.CustomerID.Valid {
		return l.claimsForUser()
	} else {
		return l.claimsForAdmin()
	}
}

func (l Login) claimsForRefresh() jwt.MapClaims {
	return jwt.MapClaims{
		"username": l.Username,
		"role": l.Role,
		"exp": time.Now().Add(REFRESH_TOKEN_DURATION).Unix(),
	}
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