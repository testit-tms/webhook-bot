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

type companyUsecases struct {
	cs companyStorage
}

var (
	ErrCompanyNotFound = errors.New("company not found")
)

func NewCompanyUsecases(cs companyStorage) *companyUsecases {
	return &companyUsecases{
		cs: cs,
	}
}

func (u *companyUsecases) GetCompanyByOwnerTelegramId(ctx context.Context, ownerId int64) (entities.Company, error) {
	const op = "usecases.GetCompanyByOwnerTelegramId"

	company, err := u.cs.GetCompanyByOwnerTelegramId(ctx, ownerId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return entities.Company{}, fmt.Errorf("%s: %w", op, ErrCompanyNotFound)
		}
		return entities.Company{}, fmt.Errorf("%s: get company by owner id: %w", op, err)
	}

	return company, nil
}
