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
	addCompany            = "INSERT INTO companies (token, owner_id, name, email) VALUES ($1, $2, $3, $4) RETURNING id, token, owner_id, name, email"
	getCompanyByOwnerId   = "SELECT c.id, c.token, c.owner_id, c.name, c.email FROM companies AS c INNER JOIN owners As o ON o.id = c.owner_id WHERE o.telegram_id=$1"
	getCompanyIdByName    = "SELECT id FROM companies WHERE name=$1"
	updateToken           = "UPDATE companies SET token=$1 WHERE id=$2"
	deleteCompany         = "DELETE FROM companies WHERE id=$1"
	deleteChatByCompanyId = "DELETE FROM chats WHERE company_id=$1"
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

// UpdateToken updates the company's token.
func (s *CompanyStorage) UpdateToken(ctx context.Context, companyId int64, token string) error {
	const op = "storage.postgres.UpdateToken"

	_, err := s.db.ExecContext(ctx, updateToken, token, companyId)
	if err != nil {
		return fmt.Errorf("%s: execute query: %w", op, err)
	}

	return nil
}

// DeleteCompany deletes a company by its ID.
func (s *CompanyStorage) DeleteCompany(ctx context.Context, companyId int64) (err error) {
	const op = "storage.postgres.DeleteCompany"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("%s: rollback transaction: %w", op, rollbackErr)
			}
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, deleteChatByCompanyId, companyId)
	if err != nil {
		return fmt.Errorf("%s: delete chat by company id: %w", op, err)
	}

	_, err = tx.ExecContext(ctx, deleteCompany, companyId)
	if err != nil {
		return fmt.Errorf("%s: delete company: %w", op, err)
	}

	return nil
}
