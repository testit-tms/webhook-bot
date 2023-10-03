package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetStartMessage(m *tgbotapi.Message) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(m.Chat.ID,
		`Hi! Welcome to the Test IT webhook bot!
	This bot helps you to get notifications from Test IT to your Telegram chat.
	To start using this bot, you need to register your company.
	To do this, use the /register command.

	Also you can see the list of available commands using the /help command.
	`)
}
