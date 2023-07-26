package commands

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/testit-tms/webhook-bot/internal/entities"
)

type companyUsesaces interface {
	GetCompaniesByOwnerId(ctx context.Context, ownerId int64) ([]entities.CompanyInfo, error)
}

type CompanyCommands struct {
	cu companyUsesaces
}

func NewCompanyCommands(cu companyUsesaces) *CompanyCommands {
	return &CompanyCommands{
		cu: cu,
	}
}

func (c *CompanyCommands) GetMyCompanies(m *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	const op = "CompanyCommands.GetMyCompanies"

	companies, err := c.cu.GetCompaniesByOwnerId(context.Background(), m.From.ID)
	if err != nil {
		return tgbotapi.NewMessage(m.Chat.ID, "Something went wrong. Lets try again"),
			fmt.Errorf("%s: get companies by owner id: %w", op, err)
	}
	msg := tgbotapi.NewMessage(m.Chat.ID, "")
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	if len(companies) == 0 {
		msg.Text = `
		*You have no companies*

		You can register new company with */register* command
		`
		return msg, nil
	}

	msg.Text = `
	*Your companies:*
	`
	for i, company := range companies {
		msg.Text += fmt.Sprintf(`
		%d\. *Name:*  _%s_ 
			 *Email:* _%s_
		`, i+1, company.Name, strings.Replace(company.Email, ".", "\\.", 1))
	}

	return msg, nil
}
