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

type CompanyStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *CompanyStorage {
	return &CompanyStorage{
		db: db,
	}
}

const (
	addCompany          = "INSERT INTO companies (token, owner_id, name, email) VALUES ($1, $2, $3, $4) RETURNING id, token, owner_id, name, email"
	getCompanyByOwnerId = "SELECT name, email FROM companies AS c INNER JOIN owners As o ON o.id = c.owner_id WHERE o.telegram_id=$1"
)

func (s *CompanyStorage) AddCompany(ctx context.Context, company entities.Company) (entities.Company, error) {
	const op = "storage.postgres.AddCompany"

	newCompany := entities.Company{}

	err := s.db.QueryRowxContext(ctx, addCompany, company.Token, company.OwnerId, company.Name, company.Email).StructScan(&newCompany)
	if err != nil {
		return newCompany, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return newCompany, nil
}

// TODO: rename to GetCompanyByOwnerTelegramId
func (s *CompanyStorage) GetCompaniesByOwnerId(ctx context.Context, ownerId int64) ([]entities.Company, error) {
	const op = "storage.postgres.GetCompanyByOwnerId"

	companies := []entities.Company{}

	err := s.db.SelectContext(ctx, &companies, getCompanyByOwnerId, ownerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return companies, storage.ErrNotFound
		}

		return companies, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return companies, nil
}
