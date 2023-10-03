package models

import "github.com/testit-tms/webhook-bot/internal/entities"

// Company represents a company that is registering for the webhook bot service.
type Company struct {
	User        User
	CompanyName string `validate:"min=1,max=100"`
	Email       string `validate:"email"`
}

// ToCompanyInfo converts a Company struct to a CompanyRegistrationInfo struct.
// It returns a new CompanyRegistrationInfo struct with the converted data.
func (c Company) ToCompanyInfo() entities.CompanyRegistrationInfo {
	return entities.CompanyRegistrationInfo{
		Name:  c.CompanyName,
		Email: c.Email,
		Owner: entities.OwnerInfo{
			TelegramID:   c.User.ID,
			TelegramName: c.User.Name,
		},
	}
}
