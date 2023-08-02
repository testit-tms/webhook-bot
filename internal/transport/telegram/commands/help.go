package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetHelpMessage(m *tgbotapi.Message) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(m.Chat.ID,
		`Available commands:
	/help - show this message
	/show - show chat ID
	/register - register new company
	/list - show registered company
	/add {chat_id} - add chat to company, for example: /add 123456789
	`)
}
