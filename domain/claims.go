package domain

import (
	"encoding/json"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

const HMAC_SAMPLE_SECRET = "hmacSampleSecret"

type Claims struct {
	CustomerID string   `json:"customer_id"`
	Accounts   []string `json:"accounts"`
	Username   string   `json:"username"`
	Expiry     int64    `json:"exp"`
	Role       string   `json:"role"`
}

func BuildClaims(MapClaims jwt.MapClaims) (*Claims, error) {
	var claims Claims
	bytes, err := json.Marshal(MapClaims)
	if err != nil {
		log.Println("Error while marshal MapClaims " + err.Error())
		return nil, err
	}
	
	err = json.Unmarshal(bytes, &claims)
	if err != nil {
		log.Println("Error while unmarshal json to Claims " + err.Error())
		return nil, err
	}

	return &claims, nil
}

func (c Claims) IsUser() bool {
	return c.Role == "user"
}

func (c Claims) IsValidAccountID(account string) bool {
	if account == "" {
		return true
	}

	for _, v := range c.Accounts {
		if account == v {
			return true
		}
	}

	return false
}

func (c Claims) IsRequestVerifiedWithTokenClaims(urlParams map[string]string) bool {
	if c.CustomerID != urlParams["id"] {
		log.Println("Not allow to access ohther customer's account")
		return false
	}

	if !c.IsValidAccountID(urlParams["account_id"]) {
		return false
	}

	return true
}