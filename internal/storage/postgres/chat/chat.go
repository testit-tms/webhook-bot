package chat

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/storage"
)

type ChatStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *ChatStorage {
	return &ChatStorage{
		db: db,
	}
}

const (
	getChatsByCompanyId    = "SELECT id, company_id, telegram_id FROM chats WHERE company_id=$1"
	getChatsByCompanyToken = "SELECT id, company_id, telegram_id FROM chats WHERE company_id=(SELECT id FROM companies WHERE token=$1)"
	addChat                = "INSERT INTO chats (company_id, telegram_id) VALUES ($1, $2) RETURNING id, company_id, telegram_id"
	deleteChatById         = "DELETE FROM chats WHERE id=$1"
	deleteChatByCompanyId  = "DELETE FROM chats WHERE company_id=$1"
)

func (s *ChatStorage) GetChatsByCompanyId(ctx context.Context, id int64) ([]entities.Chat, error) {
	const op = "storage.postgres.GetChatByCompanyId"

	chats := []entities.Chat{}

	if err := s.db.SelectContext(ctx, &chats, getChatsByCompanyId, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return chats, storage.ErrNotFound
		}

		return chats, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return chats, nil
}

func (s *ChatStorage) GetChatsByCompanyToken(ctx context.Context, t string) ([]entities.Chat, error) {
	const op = "storage.postgres.GetChatsByCompanyToken"

	chats := []entities.Chat{}

	if err := s.db.SelectContext(ctx, &chats, getChatsByCompanyToken, t); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return chats, storage.ErrNotFound
		}

		return chats, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return chats, nil
}

func (r *ChatStorage) AddChat(ctx context.Context, chat entities.Chat) (entities.Chat, error) {
	const op = "storage.postgres.AddChat"

	newChat := entities.Chat{}

	err := r.db.QueryRowxContext(ctx, addChat, chat.CompanyId, chat.TelegramId).StructScan(&newChat)
	if err != nil {
		return newChat, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return newChat, nil
}

func (r *ChatStorage) DeleteChatById(ctx context.Context, id int64) error {
	const op = "storage.postgres.DeleteChatById"

	_, err := r.db.ExecContext(ctx, deleteChatById, id)
	if err != nil {
		return fmt.Errorf("%s: execute query: %w", op, err)
	}

	return nil
}

func (r *ChatStorage) DeleteChatByCompanyId(ctx context.Context, id int) error {
	const op = "storage.postgres.DeleteChatByCompanyId"

	_, err := r.db.ExecContext(ctx, deleteChatByCompanyId, id)
	if err != nil {
		return fmt.Errorf("%s: execute query: %w", op, err)
	}

	return nil
}
