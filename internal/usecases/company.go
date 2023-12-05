package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/lib/random"
	"github.com/testit-tms/webhook-bot/internal/storage"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type companyStorage interface {
	GetCompanyByOwnerTelegramId(ctx context.Context, ownerId int64) (entities.Company, error)
	UpdateToken(ctx context.Context, companyId int64, token string) error
}

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type chatStorage interface {
	GetChatsByCompanyId(ctx context.Context, id int64) ([]entities.Chat, error)
}

type companyUsecases struct {
	cs  companyStorage
	chs chatStorage
}

var (
	// ErrCompanyNotFound is returned when a company is not found.
	ErrCompanyNotFound = errors.New("company not found")
)

// NewCompanyUsecases creates a new instance of companyUsecases.
func NewCompanyUsecases(cs companyStorage, chs chatStorage) *companyUsecases {
	return &companyUsecases{
		cs:  cs,
		chs: chs,
	}
}

// GetCompanyByOwnerTelegramId retrieves the company information associated with the given owner Telegram ID.
// It returns a CompanyInfo struct and an error. If the company is not found, it returns ErrCompanyNotFound.
func (u *companyUsecases) GetCompanyByOwnerTelegramId(ctx context.Context, ownerId int64) (entities.CompanyInfo, error) {
	const op = "usecases.GetCompanyByOwnerTelegramId"

	company, err := u.cs.GetCompanyByOwnerTelegramId(ctx, ownerId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return entities.CompanyInfo{}, fmt.Errorf("%s: %w", op, ErrCompanyNotFound)
		}
		return entities.CompanyInfo{}, fmt.Errorf("%s: get company by owner id: %w", op, err)
	}

	ci := entities.CompanyInfo{
		ID:      company.ID,
		OwnerID: company.OwnerID,
		Token:   company.Token,
		Name:    company.Name,
		Email:   company.Email,
	}

	chats, err := u.chs.GetChatsByCompanyId(ctx, company.ID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ci, nil
		}
		return ci, fmt.Errorf("%s: get chats by company id: %w", op, err)
	}

	for _, chat := range chats {
		ci.ChatIds = append(ci.ChatIds, chat.TelegramID)
	}

	return ci, nil
}

// UpdateToken updates the company token.
// It returns an error if the company is not found.
func (u *companyUsecases) UpdateToken(ctx context.Context, ownerId int64) error {
	const op = "usecases.UpdateToken"

	company, err := u.cs.GetCompanyByOwnerTelegramId(ctx, ownerId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, ErrCompanyNotFound)
		}
		return fmt.Errorf("%s: get company by owner id: %w", op, err)
	}

	token := random.NewRandomString(30)
	
	if err := u.cs.UpdateToken(ctx, company.ID, token); err != nil {
		return fmt.Errorf("%s: update token: %w", op, err)
	}

	return nil
}