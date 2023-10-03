package commands

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/usecases"
)

type companyUsesaces interface {
	GetCompanyByOwnerTelegramId(ctx context.Context, ownerId int64) (entities.CompanyInfo, error)
}

type CompanyCommands struct {
	cu companyUsesaces
}

func NewCompanyCommands(cu companyUsesaces) *CompanyCommands {
	return &CompanyCommands{
		cu: cu,
	}
}

// TODO: rename to GetMyCompany and refactor
func (c *CompanyCommands) GetMyCompanies(m *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	const op = "CompanyCommands.GetMyCompanies"

	msg := tgbotapi.NewMessage(m.Chat.ID, "")
	msg.ParseMode = tgbotapi.ModeHTML
	company, err := c.cu.GetCompanyByOwnerTelegramId(context.Background(), m.From.ID)
	if err != nil {
		if errors.Is(err, usecases.ErrCompanyNotFound) {
			msg.Text = `
			<b>You have no companies</b>
	
			You can register new company with <b>/register</b> command
			`
			return msg, nil
		}
		msg.Text = "Something went wrong. Lets try again"
		return msg, fmt.Errorf("%s: get companies by owner id: %w", op, err)
	}

	// TODO: add more better formatting and add Chats
	msg.Text = fmt.Sprintf(`
		<b>Your company:</b>
		<b>Name:</b>  <i>%s</i> 
		<b>Email:</b> <i>%s</i>
		<b>Token:</b> <i>%s</i>
		`, company.Name, company.Email, company.Token)

	if len(company.ChatIds) > 0 {
		msg.Text += "\n<b>Chats:</b>"
		for _, chatId := range company.ChatIds {
			msg.Text += fmt.Sprintf("\n%d", chatId)
		}
	}

	return msg, nil
}
