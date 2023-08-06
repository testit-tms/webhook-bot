package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetHelpMessage(m *tgbotapi.Message) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(m.Chat.ID,
		`Available commands:
	/help - show this message
	/getchatid - show chat ID
	/register - register new company
	/getcompany - show registered company
	/addchat {chat_id} - add chat to company, for example: /addchat 123456789
	/deletechat {chat_id} - delete chat from company, for example: /deletechat 123456789
	`)
}
