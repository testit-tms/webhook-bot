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
	UpdateToken(ctx context.Context, ownerId int64) error
}

// CompanyCommands represents a set of commands related to companies.
type CompanyCommands struct {
	cu companyUsesaces
}

// NewCompanyCommands creates a new instance of CompanyCommands with the provided company use cases.
func NewCompanyCommands(cu companyUsesaces) *CompanyCommands {
	return &CompanyCommands{
		cu: cu,
	}
}

// TODO: rename to GetMyCompany and refactor

// GetMyCompanies returns a Telegram message containing information about the company owned by the user who sent the message.
// If the user does not own any companies, the message will indicate that they have no companies and provide a command to register a new one.
// If an error occurs while retrieving the company information, an error message will be returned.
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

// UpdateToken updates the token of the company owned by the user who sent the message.
// If the user does not own any companies, the message will indicate that they have no companies and provide a command to register a new one.
// If an error occurs while retrieving the company information, an error message will be returned.
func (c *CompanyCommands) UpdateToken(m *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	const op = "CompanyCommands.UpdateToken"

	msg := tgbotapi.NewMessage(m.Chat.ID, "")
	msg.ParseMode = tgbotapi.ModeHTML

	err := c.cu.UpdateToken(context.Background(), m.From.ID)
	if err != nil {
		if errors.Is(err, usecases.ErrCompanyNotFound) {
			msg.Text = `
			<b>You have no companies</b>
	
			You can register new company with <b>/register</b> command
			`
			return msg, nil
		}
		msg.Text = "Something went wrong. Lets try again"
		return msg, fmt.Errorf("%s: update token: %w", op, err)
	}

	msg.Text = "Token updated successfully"
	return msg, nil
}
