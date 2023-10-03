package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetStartMessage returns a Telegram message configuration for the start command response.
// The message contains a welcome message and instructions on how to use the Test IT webhook bot.
// Parameter m is a pointer to a tgbotapi.Message struct representing the incoming message.
// The returned value is a tgbotapi.MessageConfig struct representing the message to be sent.
func GetStartMessage(m *tgbotapi.Message) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(m.Chat.ID,
		`Hi! Welcome to the Test IT webhook bot!
	This bot helps you to get notifications from Test IT to your Telegram chat.
	To start using this bot, you need to register your company.
	To do this, use the /register command.

	Also you can see the list of available commands using the /help command.
	`)
}
