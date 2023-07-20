package company

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
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
	addCompany = "INSERT INTO companies (token, owner_id, name, email) VALUES ($1, $2, $3, $4) RETURNING id, token, owner_id, name, email"
)

func (s *CompanyStorage) AddCompany(ctx context.Context, company storage.Company) (storage.Company, error) {
	const op = "storage.postgres.AddCompany"

	newCompany := storage.Company{}

	err := s.db.QueryRowxContext(ctx, addCompany, company.Token, company.OwnerId, company.Name, company.Email).StructScan(&newCompany)
	if err != nil {
		return newCompany, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return newCompany, nil
}
