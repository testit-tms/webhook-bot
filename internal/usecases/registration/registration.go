package registration

import (
	"context"
	"fmt"

	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/lib/random"
	"github.com/testit-tms/webhook-bot/internal/storage"
)

type ownerStorage interface {
	AddOwner(ctx context.Context, owner entities.Owner) (entities.Owner, error)
	GetOwnerByTelegramId(ctx context.Context, id int64) (entities.Owner, error)
}

type companyStorage interface {
	AddCompany(ctx context.Context, company entities.Company) (entities.Company, error)
}

type RegistrationUsecases struct {
	os ownerStorage
	cs companyStorage
}

func New(os ownerStorage, cs companyStorage) *RegistrationUsecases {
	return &RegistrationUsecases{
		os: os,
		cs: cs,
	}
}

// TODO: add transaction
func (r *RegistrationUsecases) RegisterCompany(ctx context.Context, c entities.CompanyRegistrationInfo) error {
	const op = "RegistrationUsecases.RegisterCompany"

	owner, err := r.os.GetOwnerByTelegramId(ctx, c.Owner.TelegramId)
	if err != nil {
		if err != storage.ErrNotFound {
			return fmt.Errorf("%s: cannot get owner by telegram id: %w", op, err)
		}

		newOwner := entities.Owner{
			TelegramId:   c.Owner.TelegramId,
			TelegramName: c.Owner.TelegramName,
		}

		owner, err = r.os.AddOwner(ctx, newOwner)
		if err != nil {
			return fmt.Errorf("%s: cannot add owner: %w", op, err)
		}
	}

	company := entities.Company{
		Name:    c.Name,
		Email:   c.Email,
		OwnerId: owner.Id,
		Token:   random.NewRandomString(30),
	}

	_, err = r.cs.AddCompany(ctx, company)
	if err != nil {
		return fmt.Errorf("%s: cannot add company: %w", op, err)
	}

	return nil
}
