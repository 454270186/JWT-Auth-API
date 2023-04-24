package domain

import (
	"database/sql"
	"errors"
	"log"
)

type AuthRepo interface {
	FindBy(username, password string) (*Login, error)
}

type AuthRepoDB struct {
	DB *sql.DB
}

func NewAuthRepo(db *sql.DB) AuthRepo {
	return AuthRepoDB{
		DB: db,
	}
}

func (ar AuthRepoDB) FindBy(username, password string) (*Login, error) {
	log.Println(username + " " + password)
	var login Login
	selectSql := `
		SELECT u.username, u.customer_id, u.role, string_agg(a.account_id::text, ',') AS account_numbers
		FROM users u
		LEFT JOIN accounts a ON u.customer_id = a.customer_id
		WHERE username = $1 AND password = $2
		GROUP BY u.username, u.customer_id, u.role;
	`

	err := ar.DB.QueryRow(selectSql, username, password).Scan(&login.Username, &login.CustomerID, &login.Role, &login.Accounts)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credential")
		}
		log.Println("Error while find by username and password " + err.Error())
		return nil, errors.New("unexpect database error")
	}

	return &login, nil
}

