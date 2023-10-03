package usecases

import (
	"context"
	"fmt"

	"github.com/testit-tms/webhook-bot/internal/entities"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type chatsStorage interface {
	AddChat(ctx context.Context, chat entities.Chat) (entities.Chat, error)
	DeleteChatById(ctx context.Context, id int64) error
	GetChatsByCompanyId(ctx context.Context, id int64) ([]entities.Chat, error)
}

type chatUsecases struct {
	cs   chatsStorage
	coms companyStorage
}

// NewChatUsecases returns a new instance of the chatUsecases struct, which provides use cases for managing chats.
// It takes in a chatsStorage interface and a companyStorage interface as parameters.
func NewChatUsecases(cs chatsStorage, coms companyStorage) *chatUsecases {
	return &chatUsecases{
		cs:   cs,
		coms: coms,
	}
}

// AddChat adds a new chat to the system.
// It takes a context and a Chat entity as input and returns the newly added Chat entity and an error (if any).
func (u *chatUsecases) AddChat(ctx context.Context, chat entities.Chat) (entities.Chat, error) {
	return u.cs.AddChat(ctx, chat)
}

// DeleteChatByTelegramId deletes a chat by its Telegram ID for a given owner ID.
// It first retrieves the company associated with the owner ID, then gets all chats
// associated with that company. If a chat with the given Telegram ID is found, it is
// deleted and the function returns nil. If no chat is found, it returns an error.
func (u *chatUsecases) DeleteChatByTelegramId(ctx context.Context, ownerId, chatId int64) error {
	const op = "usecases.DeleteChatByTelegramId"

	company, err := u.coms.GetCompanyByOwnerTelegramId(ctx, ownerId)
	if err != nil {
		return fmt.Errorf("%s: get company by owner id: %w", op, err)
	}

	chats, err := u.cs.GetChatsByCompanyId(ctx, company.Id)
	if err != nil {
		return fmt.Errorf("%s: get chats by company id: %w", op, err)
	}

	for _, chat := range chats {
		if chat.TelegramId == chatId {
			if err := u.cs.DeleteChatById(ctx, chat.Id); err != nil {
				return fmt.Errorf("%s: delete chat by id: %w", op, err)
			}
			return nil
		}
	}

	return fmt.Errorf("%s: chat not found", op)
}
