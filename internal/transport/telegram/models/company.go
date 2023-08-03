package models

import "github.com/testit-tms/webhook-bot/internal/entities"

type Company struct {
	User        User
	CompanyName string `validate:"min=1,max=100"`
	Email       string `validate:"email"`
}

func (c Company) ToCompanyInfo() entities.CompanyRegistrationInfo {
	return entities.CompanyRegistrationInfo{
		Name:  c.CompanyName,
		Email: c.Email,
		Owner: entities.OwnerInfo{
			TelegramId:   c.User.ID,
			TelegramName: c.User.Name,
		},
	}
}
