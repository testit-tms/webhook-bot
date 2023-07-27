package commands

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/testit-tms/webhook-bot/internal/entities"
)

type chatUsecases interface {
	AddChat(ctx context.Context, chat entities.Chat) (entities.Chat, error)
}

type chatCommands struct {
	cu    chatUsecases
	compu companyUsesaces
}

func NewChatCommands(cu chatUsecases, compu companyUsesaces) *chatCommands {
	return &chatCommands{
		cu:    cu,
		compu: compu,
	}
}

func (c *chatCommands) AddChat(m *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	const op = "chatCommands.AddChat"

	args := m.CommandArguments()
	chatID, err := strconv.Atoi(args)
	if err != nil {
		return tgbotapi.NewMessage(m.Chat.ID, "Wrong chat id"),
			fmt.Errorf("%s: convert chat id: %w", op, err)
	}

	company, err := c.compu.GetCompanyByOwnerTelegramId(context.Background(), m.From.ID)
	if err != nil {
		return tgbotapi.NewMessage(m.Chat.ID, "Something went wrong. Lets try again"),
			fmt.Errorf("%s: get company by owner id: %w", op, err)
	}

	_, err = c.cu.AddChat(context.Background(), entities.Chat{
		CompanyId:  company.Id,
		TelegramId: int64(chatID),
	})
	if err != nil {
		return tgbotapi.NewMessage(m.Chat.ID, "Something went wrong. Lets try again"),
			fmt.Errorf("%s: add chat: %w", op, err)
	}

	return tgbotapi.NewMessage(m.Chat.ID, "Chat added"), nil
}
