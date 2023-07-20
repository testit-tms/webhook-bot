package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/testit-tms/webhook-bot/internal/storage"
)

type OwnerStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *OwnerStorage {
	return &OwnerStorage{
		db: db,
	}
}

const (
	getOwnerById = "SELECT id, telegram_id, telegram_name FROM owners WHERE id=$1"
)

func (s *OwnerStorage) GetOwnerById(ctx context.Context, id int) (storage.Owner, error) {
	const op = "storage.postgres.GetOwnerById"

	owner := storage.Owner{}

	if err := s.db.GetContext(ctx, &owner, getOwnerById, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return owner, storage.ErrNotFound
		}

		return owner, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return owner, nil
}
