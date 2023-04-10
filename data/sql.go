package data

import (
	"database/sql"
	"os"
)

func OpenDB() (*sql.DB, error) {
	dsn := os.Getenv("DSN")

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
