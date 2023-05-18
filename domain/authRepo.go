package domain

import (
	"database/sql"
	"errors"
	"jwtAuth/dto"
	"log"
)

type AuthRepo interface {
	FindBy(username, password string) (*Login, error)
	FindByName(username string) (*Login, error)
	StoreRefreshToken(userToken AuthToken) error
	IsRefreshTokenExist(refreshReq dto.RefreshRequest) error
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

func (ar AuthRepoDB) FindByName(username string) (*Login, error) {
	var login Login
	selectSql := `
		SELECT u.username, u.customer_id, u.role, string_agg(a.account_id::text, ',') AS account_numbers
		FROM users u
		LEFT JOIN accounts a ON u.customer_id = a.customer_id
		WHERE username = $1
		GROUP BY u.username, u.customer_id, u.role;
	`

	err := ar.DB.QueryRow(selectSql, username).Scan(&login.Username, &login.CustomerID, &login.Role, &login.Accounts)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credential")
		}
		log.Println("Error while find by username " + err.Error())
		return nil, errors.New("unexpect database error")
	}

	return &login, nil
}

func (ar AuthRepoDB) StoreRefreshToken(userToken AuthToken) error {
	log.Println(userToken.Token + "----" + userToken.RefreshToken)

	insertSql := `
		INSERT INTO refresh_token (refresh_token)
		VALUES ($1)
	`

	_, err := ar.DB.Exec(insertSql, userToken.RefreshToken)
	if err != nil {
		log.Println("Error while store a refresh token")
		return errors.New("unexpect database error")
	}

	return nil
}

func (ar AuthRepoDB) IsRefreshTokenExist(refreshReq dto.RefreshRequest) error {
	selectSql := `
		SELECT refresh_token FROM refresh_token 
		WHERE refresh_token = $1
	`

	var reToken string
	err := ar.DB.QueryRow(selectSql, refreshReq.RefreshToken).Scan(&reToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		log.Println("Error while find refresh token")
		return errors.New("unexpect database error")
	}

	return nil
}

