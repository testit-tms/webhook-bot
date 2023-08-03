package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/storage"
	"golang.org/x/exp/slog"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type chatGeter interface {
	GetChatsByCompanyToken(ctx context.Context, t string) ([]entities.Chat, error)
}

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type botSender interface {
	SendMessage(ctx context.Context, msg entities.Message) error
}

type sendMessageUsacases struct {
	logger *slog.Logger
	cg     chatGeter
	bs     botSender
}

var (
	ErrChatsNotFound = errors.New("chats not found")
	ErrChatsNotAllow = errors.New("chats not allowed")
	ErrCanNotSend    = errors.New("can not send message")
)

func NewSendMessageUsecases(logger *slog.Logger, cg chatGeter, bs botSender) *sendMessageUsacases {
	return &sendMessageUsacases{
		logger: logger,
		cg:     cg,
		bs:     bs,
	}
}

func (u *sendMessageUsacases) SendMessage(ctx context.Context, msg entities.Message) error {
	const op = "usecases.SendMessage"
	logger := u.logger.With(slog.String("operation", op))

	chats, err := u.cg.GetChatsByCompanyToken(ctx, msg.Token)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			logger.Debug("chats not found")
			return fmt.Errorf("%s: chats not found: %w", op, ErrChatsNotFound)
		}
		logger.Error("get chats by company token", "error", err)
		return fmt.Errorf("%s: get chats by company token: %w", op, ErrChatsNotFound)
	}

	if len(msg.ChatIds) == 0 {
		for _, c := range chats {
			msg.ChatIds = append(msg.ChatIds, c.TelegramId)
		}
		err := u.bs.SendMessage(ctx, msg)
		if err != nil {
			logger.Error("can not send message", "error", err)
			return fmt.Errorf("%s: can not send message: %w", op, ErrCanNotSend)
		}

		return nil
	}

	allowedChats := make([]int64, 0, len(msg.ChatIds))
	for _, chat := range msg.ChatIds {
		for _, c := range chats {
			if c.TelegramId == chat {
				allowedChats = append(allowedChats, c.TelegramId)
			}
		}
	}

	if len(allowedChats) == 0 {
		logger.Debug("chats not allowed")
		return fmt.Errorf("%s: chats not allowed: %w", op, ErrChatsNotAllow)
	}

	msg.ChatIds = allowedChats
	err = u.bs.SendMessage(ctx, msg)
	if err != nil {
		logger.Error("can not send message", "error", err)
		return fmt.Errorf("%s: can not send message: %w", op, ErrCanNotSend)
	}

	return nil
}
