package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	company, err := c.cu.GetCompanyByOwnerTelegramId(context.Background(), m.From.ID)
	if err != nil {
		if errors.Is(err, usecases.ErrCompanyNotFound) {
			msg.Text = `
			*You have no companies*
	
			You can register new company with */register* command
			`
			return msg, nil
		}
		msg.Text = "Something went wrong. Lets try again"
		return msg, fmt.Errorf("%s: get companies by owner id: %w", op, err)
	}

	// TODO: add more better formatting and add Chats
	msg.Text = fmt.Sprintf(`
		*Your company:*
		*Name:*  _%s_ 
		*Email:* _%s_
		*Token:* _%s_
		`, replaceSpecialCharacters(company.Name), strings.Replace(company.Email, ".", "\\.", -1), company.Token)

	if len(company.ChatIds) > 0 {
		msg.Text += "\n*Chats:*"
		for _, chatId := range company.ChatIds {
			msg.Text += strings.Replace(fmt.Sprintf("\n%d", chatId), "-", "\\-", 1)
		}
	}

	return msg, nil
}

func replaceSpecialCharacters(str string) string {
	for _, char := range []string{"-", "[", "]", "(", ")", "~", "`", ">", "#", "+", "=", "|", "{", "}", ".", "!"} {
		str = strings.Replace(str, char, "\\"+char, -1)
	}
	return str
}
