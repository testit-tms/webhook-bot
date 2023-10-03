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
	DeleteChatByTelegramId(ctx context.Context, ownerId, chatId int64) error
}

type chatCommands struct {
	cu    chatUsecases
	compu companyUsesaces
}

// NewChatCommands returns a new instance of chatCommands with the provided chatUsecases and companyUsesaces.
func NewChatCommands(cu chatUsecases, compu companyUsesaces) *chatCommands {
	return &chatCommands{
		cu:    cu,
		compu: compu,
	}
}

// AddChat adds a new chat to the company with the owner's Telegram ID.
// It takes a Telegram message as input and extracts the chat ID from the command arguments.
// If the chat ID is not a valid integer, it returns an error.
// It then retrieves the company associated with the owner's Telegram ID and adds the chat to the company.
// If there is an error while adding the chat, it returns an error.
// Otherwise, it returns a success message.
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
		CompanyID:  company.ID,
		TelegramID: int64(chatID),
	})
	if err != nil {
		return tgbotapi.NewMessage(m.Chat.ID, "Something went wrong. Lets try again"),
			fmt.Errorf("%s: add chat: %w", op, err)
	}

	return tgbotapi.NewMessage(m.Chat.ID, "Chat added"), nil
}

// DeleteChat deletes a chat by its ID. It takes a Telegram message as input and extracts the chat ID from the command arguments.
// It then calls the DeleteChatByTelegramId method of the ChatUseCase to delete the chat from the database.
// If the chat ID is not a valid integer, it returns an error and a message to the user.
// If there is an error deleting the chat, it returns an error and a message to the user.
// Otherwise, it returns a success message to the user.
func (c *chatCommands) DeleteChat(m *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	const op = "chatCommands.DeleteChat"

	args := m.CommandArguments()
	chatID, err := strconv.Atoi(args)
	if err != nil {
		return tgbotapi.NewMessage(m.Chat.ID, "Wrong chat id"),
			fmt.Errorf("%s: convert chat id: %w", op, err)
	}

	if err := c.cu.DeleteChatByTelegramId(context.Background(), m.From.ID, int64(chatID)); err != nil {
		return tgbotapi.NewMessage(m.Chat.ID, "Something went wrong. Lets try again"),
			fmt.Errorf("%s: delete chat by telegram id: %w", op, err)
	}

	return tgbotapi.NewMessage(m.Chat.ID, "Chat deleted"), nil
}
