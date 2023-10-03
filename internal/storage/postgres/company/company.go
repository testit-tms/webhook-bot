package company

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/storage"
)

// CompanyStorage is a storage implementation for companies using PostgreSQL.
type CompanyStorage struct {
	db *sqlx.DB
}

// New creates a new instance of CompanyStorage with the given database connection.
func New(db *sqlx.DB) *CompanyStorage {
	return &CompanyStorage{
		db: db,
	}
}

const (
	addCompany          = "INSERT INTO companies (token, owner_id, name, email) VALUES ($1, $2, $3, $4) RETURNING id, token, owner_id, name, email"
	getCompanyByOwnerId = "SELECT c.id, c.token, c.owner_id, c.name, c.email FROM companies AS c INNER JOIN owners As o ON o.id = c.owner_id WHERE o.telegram_id=$1"
	getCompanyIdByName  = "SELECT id FROM companies WHERE name=$1"
)

// AddCompany adds a new company to the database and returns the newly created company.
func (s *CompanyStorage) AddCompany(ctx context.Context, company entities.Company) (entities.Company, error) {
	const op = "storage.postgres.AddCompany"

	newCompany := entities.Company{}

	err := s.db.QueryRowxContext(ctx, addCompany, company.Token, company.OwnerID, company.Name, company.Email).StructScan(&newCompany)
	if err != nil {
		return newCompany, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return newCompany, nil
}

// GetCompanyByOwnerTelegramId retrieves a company by the owner's Telegram ID.
// If the company is not found, ErrNotFound is returned.
func (s *CompanyStorage) GetCompanyByOwnerTelegramId(ctx context.Context, ownerId int64) (entities.Company, error) {
	const op = "storage.postgres.GetCompanyByOwnerId"

	company := entities.Company{}

	err := s.db.GetContext(ctx, &company, getCompanyByOwnerId, ownerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return company, storage.ErrNotFound
		}

		return company, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return company, nil
}
