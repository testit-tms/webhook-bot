package usecases

import (
	"context"
	"fmt"

	"github.com/testit-tms/webhook-bot/internal/entities"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type companyStorage interface {
	GetCompaniesByOwnerId(ctx context.Context, ownerId int64) ([]entities.Company, error)
}

type companyUsecases struct {
	cs companyStorage
}

func NewCompanyUsecases(cs companyStorage) *companyUsecases {
	return &companyUsecases{
		cs: cs,
	}
}

func (u *companyUsecases) GetCompaniesByOwnerId(ctx context.Context, ownerId int64) ([]entities.CompanyInfo, error) {
	const op = "usecases.GetCompaniesByOwnerId"

	companies, err := u.cs.GetCompaniesByOwnerId(ctx, ownerId)
	if err != nil {
		return nil, fmt.Errorf("%s: get companies by owner id: %w", op, err)
	}

	companiesInfo := []entities.CompanyInfo{}
	for _, company := range companies {
		companiesInfo = append(companiesInfo, entities.CompanyInfo{
			Name:  company.Name,
			Email: company.Email,
		})
	}

	return companiesInfo, nil
}
