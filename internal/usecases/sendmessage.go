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
	// ErrChatsNotFound is returned when chats are not found.
	ErrChatsNotFound = errors.New("chats not found")
	// ErrChatsNotAllow is returned when chats are not allowed.
	ErrChatsNotAllow = errors.New("chats not allowed")
	// ErrCanNotSend is returned when a message cannot be sent.
	ErrCanNotSend    = errors.New("can not send message")
)

// NewSendMessageUsecases creates a new instance of sendMessageUsacases with the provided dependencies.
func NewSendMessageUsecases(logger *slog.Logger, cg chatGeter, bs botSender) *sendMessageUsacases {
	return &sendMessageUsacases{
		logger: logger,
		cg:     cg,
		bs:     bs,
	}
}

// SendMessage sends a message to the specified chats. If no chat IDs are provided, the message is sent to all chats associated with the company token.
// If chat IDs are provided, the message is only sent to the chats that are associated with the company token and have a matching chat ID.
// Returns an error if the chats are not found, not allowed, or if the message cannot be sent.
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
			msg.ChatIds = append(msg.ChatIds, c.TelegramID)
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
			if c.TelegramID == chat {
				allowedChats = append(allowedChats, c.TelegramID)
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
