package domain

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

var dbUser string
var dbPassword string
var dbHost string
var dbPort string
var dbName string

func init() {
	dbUser = os.Getenv("db_user")
	dbPassword = os.Getenv("db_password")
	dbHost = os.Getenv("db_host")
	dbPort = os.Getenv("db_port")
	dbName = os.Getenv("db_name")
}

func NewDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifeTime)

	return db, nil
}