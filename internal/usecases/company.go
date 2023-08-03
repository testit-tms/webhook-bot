package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/storage"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type companyStorage interface {
	GetCompanyByOwnerTelegramId(ctx context.Context, ownerId int64) (entities.Company, error)
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
	ErrCompanyNotFound = errors.New("company not found")
)

func NewCompanyUsecases(cs companyStorage, chs chatStorage) *companyUsecases {
	return &companyUsecases{
		cs:  cs,
		chs: chs,
	}
}

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
		Id:      company.Id,
		OwnerId: company.OwnerId,
		Token:   company.Token,
		Name:    company.Name,
		Email:   company.Email,
	}

	chats, err := u.chs.GetChatsByCompanyId(ctx, company.Id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ci, nil
		}
		return ci, fmt.Errorf("%s: get chats by company id: %w", op, err)
	}

	for _, chat := range chats {
		ci.ChatIds = append(ci.ChatIds, chat.TelegramId)
	}

	return ci, nil
}
