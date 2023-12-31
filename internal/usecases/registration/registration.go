package registration

import (
	"context"
	"errors"
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
	GetCompanyByOwnerTelegramId(ctx context.Context, ownerId int64) (entities.Company, error)
}

var (
	// ErrCompanyAlreadyExists is returned when a company already exists.
	ErrCompanyAlreadyExists = fmt.Errorf("company already exists")
)

// TODO: move to usesases package and write tests

// RegistrationUsecases is a usecase implementation for registration.
type RegistrationUsecases struct {
	os ownerStorage
	cs companyStorage
}

// New creates a new instance of RegistrationUsecases.
func New(os ownerStorage, cs companyStorage) *RegistrationUsecases {
	return &RegistrationUsecases{
		os: os,
		cs: cs,
	}
}

// TODO: move to usesases package and write tests

// CheckCompanyExists checks if a company exists for the given owner ID.
func (r *RegistrationUsecases) CheckCompanyExists(ctx context.Context, ownerId int64) (bool, error) {
	const op = "RegistrationUsecases.CheckCompanyExists"

	_, err := r.cs.GetCompanyByOwnerTelegramId(ctx, ownerId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("%s: get company by owner id: %w", op, err)
	}

	return true, nil
}

// TODO: add transaction and tests

// RegisterCompany registers a new company with the given registration information.
func (r *RegistrationUsecases) RegisterCompany(ctx context.Context, c entities.CompanyRegistrationInfo) error {
	const op = "RegistrationUsecases.RegisterCompany"

	owner, err := r.os.GetOwnerByTelegramId(ctx, c.Owner.TelegramID)
	if err != nil {
		if err != storage.ErrNotFound {
			return fmt.Errorf("%s: cannot get owner by telegram id: %w", op, err)
		}

		newOwner := entities.Owner{
			TelegramID:   c.Owner.TelegramID,
			TelegramName: c.Owner.TelegramName,
		}

		owner, err = r.os.AddOwner(ctx, newOwner)
		if err != nil {
			return fmt.Errorf("%s: cannot add owner: %w", op, err)
		}
	}

	_, err = r.cs.GetCompanyByOwnerTelegramId(ctx, c.Owner.TelegramID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			company := entities.Company{
				Name:    c.Name,
				Email:   c.Email,
				OwnerID: owner.ID,
				Token:   random.NewRandomString(30),
			}

			_, err = r.cs.AddCompany(ctx, company)
			if err != nil {
				return fmt.Errorf("%s: cannot add company: %w", op, err)
			}
			return nil
		}

		return fmt.Errorf("%s: get company by owner id: %w", op, err)
	}

	return ErrCompanyAlreadyExists
}
