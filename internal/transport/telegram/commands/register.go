package commands

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/lib/logger/sl"
	"github.com/testit-tms/webhook-bot/internal/transport/telegram/models"
	"golang.org/x/exp/slog"
)

type regUsecase interface {
	RegisterCompany(ctx context.Context, c entities.CompanyRegistrationInfo) error
}

type Registrator struct {
	logger  *slog.Logger
	waitMap map[int64]models.Company
	u       regUsecase
}

func NewRegistrator(logger *slog.Logger, u regUsecase) *Registrator {
	return &Registrator{
		logger:  logger,
		waitMap: make(map[int64]models.Company),
		u:       u,
	}
}

func (r *Registrator) registerCompany(chatID int64, companyName string, userID int64, userName string) {
	r.waitMap[chatID] = models.Company{
		CompanyName: companyName,
		User: models.User{
			ID:   userID,
			Name: userName,
		},
	}
}

func (r *Registrator) registerEmail(chatID int64, email string) {
	if c, ok := r.waitMap[chatID]; ok {
		r.waitMap[chatID] = models.Company{
			CompanyName: c.CompanyName,
			User:        c.User,
			Email:       email,
		}
	}
}

func (r *Registrator) Action(m *tgbotapi.Message, step int) (tgbotapi.MessageConfig, int) {
	const op = "Registrator.Action"
	logger := r.logger.With(
		slog.String("op", op),
	)

	switch step {
	case 1:
		r.registerCompany(m.Chat.ID, m.Text, m.From.ID, m.From.UserName)
		msg := tgbotapi.NewMessage(m.Chat.ID, "Enter email:")
		msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
		return msg, 2
	case 2:
		r.registerEmail(m.Chat.ID, m.Text)
		company, ok := r.waitMap[m.Chat.ID]
		if !ok {
			logger.Error("company not found", slog.Int64("chat_id", m.Chat.ID))
			return tgbotapi.NewMessage(m.Chat.ID, "Something went wrong. Lets try again"), 0
		}
		err := validator.New().Struct(company)
		if err != nil {
			// TODO: return formated validation errors, not standart ValidationErrors
			validationErrors := err.(validator.ValidationErrors)
			logger.Error("validation error", sl.Err(err))
			return tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("validation error: %s", validationErrors)), 0
		}
		err = r.u.RegisterCompany(context.Background(), company.ToCompanyInfo())
		if err != nil {
			logger.Error("register company", sl.Err(err))
			return tgbotapi.NewMessage(m.Chat.ID, "Something went wrong. Lets try again"), 0
		}
		msg := tgbotapi.NewMessage(m.Chat.ID, "You are registered!")
		return msg, 0
	default:
		logger.Error("unknown step", slog.Int("step", step))
		return tgbotapi.NewMessage(m.Chat.ID, "Something went wrong. Lets try again"), 0
	}
}

func (r *Registrator) GetFirstMessage(chatID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, "Enter company name:")
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
	return msg
}
