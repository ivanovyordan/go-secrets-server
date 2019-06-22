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
	var err error

	dsn, err := getDSN()
	if err != nil {
		return err
	}

	Connection, err = sql.Open("postgres", dsn)

	if err != nil {
		return err
	}

	return nil
}

func getDSN() (string, error) {
	environment, _ := os.LookupEnv("ENVIRONMENT")
	if environment != "production" && godotenv.Load() != nil {
		return "", errors.New("Error loading .env file")
	}

	dsn, exists := os.LookupEnv("DATABASE_URL")
	if exists {
		return dsn, nil
	}

	dsn = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DATABASE"),
	)

	return dsn, nil
}
