package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Initialize(host string, port int64, user, password, dbname string) (*sqlx.DB, error) {
	const op = "database.Initialize"
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname))
	if err != nil {
		return nil, fmt.Errorf("%s: cannot connect to database: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: cannot ping database: %w", op, err)
	}

	return db, nil
}
