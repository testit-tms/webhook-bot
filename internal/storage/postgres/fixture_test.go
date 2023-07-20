package postgres

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type Fixture struct {
	mock sqlmock.Sqlmock
	db   *sqlx.DB
	con  *sql.DB
}

func NewFixture(t *testing.T) *Fixture {
	mockDB, mock, _ := sqlmock.New()

	db := sqlx.NewDb(mockDB, "sqlmock")

	return &Fixture{mock: mock, db: db, con: mockDB}
}

func (f *Fixture) Teardown() {
	f.con.Close()
}
