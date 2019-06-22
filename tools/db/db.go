package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Connection *sql.DB

func Connect() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DATABASE"),
	)

	Connection, err = sql.Open("postgres", dsn)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
