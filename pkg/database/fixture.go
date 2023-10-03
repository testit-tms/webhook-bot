package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

// Fixture represents a test fixture for database testing.
type Fixture struct {
	Mock sqlmock.Sqlmock
	DB   *sqlx.DB
	con  *sql.DB
}

// NewFixture returns a new instance of Fixture for testing purposes.
// It creates a new sqlmock database connection and returns a Fixture instance
// with the sqlmock DB and connection.
func NewFixture(t *testing.T) *Fixture {
	mockDB, mock, _ := sqlmock.New()

	db := sqlx.NewDb(mockDB, "sqlmock")

	return &Fixture{Mock: mock, DB: db, con: mockDB}
}

// Teardown closes the connection to the mock database.
func (f *Fixture) Teardown() {
	f.con.Close()
}
