package owner

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/testit-tms/webhook-bot/internal/entities"
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
	getOwnerById         = "SELECT id, telegram_id, telegram_name FROM owners WHERE id=$1"
	getOwnerByTelegramId = "SELECT id, telegram_id, telegram_name FROM owners WHERE telegram_id=$1"
	addOwner             = "INSERT INTO owners (telegram_id, telegram_name) VALUES ($1, $2) RETURNING id, telegram_id, telegram_name"
	deleteOwnerById      = "DELETE FROM owners WHERE id=$1"
)

func (s *OwnerStorage) GetOwnerById(ctx context.Context, id int64) (entities.Owner, error) {
	const op = "storage.postgres.GetOwnerById"

	owner := entities.Owner{}

	if err := s.db.GetContext(ctx, &owner, getOwnerById, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return owner, storage.ErrNotFound
		}

		return owner, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return owner, nil
}

func (s *OwnerStorage) GetOwnerByTelegramId(ctx context.Context, id int64) (entities.Owner, error) {
	const op = "storage.postgres.GetOwnerByTelegramId"

	owner := entities.Owner{}

	if err := s.db.GetContext(ctx, &owner, getOwnerByTelegramId, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return owner, storage.ErrNotFound
		}

		return owner, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return owner, nil
}

func (r *OwnerStorage) AddOwner(ctx context.Context, owner entities.Owner) (entities.Owner, error) {
	const op = "storage.postgres.AddOwner"

	newOwner := entities.Owner{}

	err := r.db.QueryRowxContext(ctx, addOwner, owner.TelegramId, owner.TelegramName).StructScan(&newOwner)
	if err != nil {
		return newOwner, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return newOwner, nil
}

func (r *OwnerStorage) DeleteOwnerById(ctx context.Context, id int64) error {
	const op = "storage.postgres.DeleteOwnerById"

	_, err := r.db.ExecContext(ctx, deleteOwnerById, id)
	if err != nil {
		return fmt.Errorf("%s: execute query: %w", op, err)
	}

	return nil
}
