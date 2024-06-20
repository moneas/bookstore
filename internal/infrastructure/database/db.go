package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDB(dsn string) error {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}

	DB = db
	fmt.Println("Database connection established")

	return nil
}
