package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type Fixture struct {
	Mock sqlmock.Sqlmock
	DB   *sqlx.DB
	con  *sql.DB
}

func NewFixture(t *testing.T) *Fixture {
	mockDB, mock, _ := sqlmock.New()

	db := sqlx.NewDb(mockDB, "sqlmock")

	return &Fixture{Mock: mock, DB: db, con: mockDB}
}

func (f *Fixture) Teardown() {
	f.con.Close()
}
