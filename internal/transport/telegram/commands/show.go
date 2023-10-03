package commands

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetChatId returns a MessageConfig that contains the chat ID of the given message.
// The chat ID is included in the message text.
func GetChatId(m *tgbotapi.Message) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("Chat ID: %d", m.Chat.ID))
}
