package usecases

import (
	"context"

	"github.com/testit-tms/webhook-bot/internal/entities"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type chatsStorage interface {
	AddChat(ctx context.Context, chat entities.Chat) (entities.Chat, error)
}

type chatUsecases struct {
	cs chatsStorage
}

func NewChatUsecases(cs chatsStorage) *chatUsecases {
	return &chatUsecases{
		cs: cs,
	}
}

// TODO: add test
func (u *chatUsecases) AddChat(ctx context.Context, chat entities.Chat) (entities.Chat, error) {
	return u.cs.AddChat(ctx, chat)
}
