package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/lib/logger/sl"
	"github.com/testit-tms/webhook-bot/internal/transport/telegram/commands"
	"golang.org/x/exp/slog"
)

const (
	rigesterCommand = "register"
	getChatID       = "get_chat_id"
	getMyCompanies  = "get_my_companies"
	addChat         = "add_chat"
)

type registrator interface {
	Action(m *tgbotapi.Message, step int) (tgbotapi.MessageConfig, int)
	GetFirstMessage(chatID int64) tgbotapi.MessageConfig
}

type companyCommands interface {
	GetMyCompanies(m *tgbotapi.Message) (tgbotapi.MessageConfig, error)
}

type chatCommands interface {
	AddChat(m *tgbotapi.Message) (tgbotapi.MessageConfig, error)
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

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case rigesterCommand:
			msg = b.registrator.GetFirstMessage(update.Message.Chat.ID)
			b.waitConversation[update.Message.Chat.ID] = Conversation{
				typeOfConversation: rigesterType,
				step:               1,
			}
		case getChatID:
			msg = commands.GetChatId(update.Message)
		case getMyCompanies:
			msg, err := b.cc.GetMyCompanies(update.Message)
			if err != nil {
				b.logger.Error("cannot get company", sl.Err(err))
			}
			b.sendMessage(msg)
			continue
		case addChat:
			msg, err := b.chc.AddChat(update.Message)
			if err != nil {
				b.logger.Error("cannot add chat", sl.Err(err))
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

func (b *TelegramBot) SendMessage(channelID int64, message entities.Message) error {
	msg := tgbotapi.NewMessage(channelID, message.Text)

	if message.ParseMode != entities.Undefined {
		msg.ParseMode = string(message.ParseMode)
	}

	if _, err := b.bot.Send(msg); err != nil {
		return err
	}

	return nil
}
