package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetHelpMessage returns a Telegram message configuration containing a list of available commands.
// The message includes the command names and their descriptions.
// The message is sent to the chat ID specified in the input message.
func GetHelpMessage(m *tgbotapi.Message) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(m.Chat.ID,
		`Available commands:
	/help - show this message
	/getchatid - show chat ID
	/register - register new company
	/getcompany - show registered company
	/updatetoken - update company token
	/addchat {chat_id} - add chat to company, for example: /addchat 123456789
	/deletechat {chat_id} - delete chat from company, for example: /deletechat 123456789
	`)
}
