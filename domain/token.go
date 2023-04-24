package domain

const HMAC_SAMPLE_SECRET = "hmacSampleSecret"

type Claims struct {
	CustomerID string   `json:"customer_id"`
	Accounts   []string `json:"accounts"`
	Username   string   `json:"username"`
	Expiry     int64    `json:"exp"`
	Role       string   `json:"role"`
}
