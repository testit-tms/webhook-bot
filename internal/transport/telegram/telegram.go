package telegram

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/lib/logger/sl"
	"github.com/testit-tms/webhook-bot/internal/transport/telegram/commands"
	"golang.org/x/exp/slog"
)

const (
	rigesterCommand   = "register"
	getChatIdCommand  = "getchatid"
	getCompanyCommand = "getcompany"
	addChatCommand    = "addchat"
	helpCommand       = "help"
	deleteChatCommand = "deletechat"
	startCommand      = "start"
)

type registrator interface {
	Action(m *tgbotapi.Message, step int) (tgbotapi.MessageConfig, int)
	GetFirstMessage(m *tgbotapi.Message) tgbotapi.MessageConfig
}

type companyCommands interface {
	GetMyCompanies(m *tgbotapi.Message) (tgbotapi.MessageConfig, error)
}

type chatCommands interface {
	AddChat(m *tgbotapi.Message) (tgbotapi.MessageConfig, error)
	DeleteChat(m *tgbotapi.Message) (tgbotapi.MessageConfig, error)
}

type TelegramBot struct {
	logger           *slog.Logger
	bot              *tgbotapi.BotAPI
	waitConversation map[int64]Conversation
	registrator      registrator
	cc               companyCommands
	chc              chatCommands
}

// New creates a new TelegramBot instance
func New(logger *slog.Logger, token string, r registrator, cc companyCommands, chc chatCommands) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
			logger:           logger,
			bot:              bot,
			waitConversation: make(map[int64]Conversation),
			registrator:      r,
			cc:               cc,
			chc:              chc,
		},
		nil
}

// Run starts the bot
func (b *TelegramBot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			if conversation, ok := b.waitConversation[update.Message.Chat.ID]; ok {
				switch conversation.typeOfConversation {
				case rigesterType:
					msg, step := b.registrator.Action(update.Message, conversation.step)
					if step == 0 {
						delete(b.waitConversation, update.Message.Chat.ID)
					} else {
						b.waitConversation[update.Message.Chat.ID] = Conversation{
							typeOfConversation: conversation.typeOfConversation,
							step:               step,
						}
					}
					b.sendMessage(msg)
				}
			}
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case helpCommand:
			msg = commands.GetHelpMessage(update.Message)
		case startCommand:
			msg = commands.GetStartMessage(update.Message)
		case rigesterCommand:
			msg = b.registrator.GetFirstMessage(update.Message)
			b.waitConversation[update.Message.Chat.ID] = Conversation{
				typeOfConversation: rigesterType,
				step:               1,
			}
		case getChatIdCommand:
			msg = commands.GetChatId(update.Message)
		case getCompanyCommand:
			msg, err := b.cc.GetMyCompanies(update.Message)
			if err != nil {
				b.logger.Error("cannot get company", sl.Err(err))
			}
			b.sendMessage(msg)
			continue
		case addChatCommand:
			msg, err := b.chc.AddChat(update.Message)
			if err != nil {
				b.logger.Error("cannot add chat", sl.Err(err))
			}
			b.sendMessage(msg)
			continue
		case deleteChatCommand:
			msg, err := b.chc.DeleteChat(update.Message)
			if err != nil {
				b.logger.Error("cannot delete chat", sl.Err(err))
			}
			b.sendMessage(msg)
			continue
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := b.bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func (b *TelegramBot) sendMessage(m tgbotapi.MessageConfig) {
	const op = "telegram.sendMessage"

	if _, err := b.bot.Send(m); err != nil {
		b.logger.Error("cannot send message", sl.Err(err), slog.String("op", op))
	}
}

func (b *TelegramBot) SendMessage(ctx context.Context, msg entities.Message) error {
	const op = "telegram.SendMessage"
	for _, chatID := range msg.ChatIds {
		newMessage := tgbotapi.NewMessage(chatID, msg.Text)

		if msg.ParseMode != entities.Undefined {
			newMessage.ParseMode = string(msg.ParseMode)
		}

		if _, err := b.bot.Send(newMessage); err != nil {
			return fmt.Errorf("%s :cannot send message to chat %d: %w", op, chatID, err)
		}
	}
	return nil
}
