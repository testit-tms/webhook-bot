package database

import (
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const (
	driverName = "postgres"
)

func Initialize(host string, port int64, user, password, dbname string) (*sqlx.DB, error) {
	const op = "database.Initialize"
	db, err := sqlx.Connect(driverName, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname))
	if err != nil {
		return nil, fmt.Errorf("%s: cannot connect to database: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: cannot ping database: %w", op, err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(driverName); err != nil {
		return nil, fmt.Errorf("%s: cannot set goose dialect: %w", op, err)
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return nil, fmt.Errorf("%s: cannot run goose migrations: %w", op, err)
	}

	return db, nil
}
