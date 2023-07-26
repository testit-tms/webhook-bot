package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/storage"
	"github.com/testit-tms/webhook-bot/internal/usecases/mocks"
	"go.uber.org/mock/gomock"
)

func Test_companyUsecases_GetCompaniesByOwnerId(t *testing.T) {
	tests := []struct {
		name           string
		ownerId        int64
		want           []entities.CompanyInfo
		mockEntities   []entities.Company
		mockError      error
		wantErr        bool
		wantErrMessage string
	}{
		{
			name:    "success",
			ownerId: 1,
			want: []entities.CompanyInfo{
				{
					Name:  "Yandex",
					Email: "info@ya.ru",
				},
				{
					Name:  "Google",
					Email: "info@google.com",
				},
			},
			mockEntities: []entities.Company{
				{
					Name:  "Yandex",
					Email: "info@ya.ru",
				},
				{
					Name:  "Google",
					Email: "info@google.com",
				},
			},
			mockError: nil,
			wantErr:   false,
		},
		{
			name:           "with ErrNotFound",
			ownerId:        1,
			want:           nil,
			mockEntities:   nil,
			mockError:      storage.ErrNotFound,
			wantErr:        true,
			wantErrMessage: "usecases.GetCompaniesByOwnerId: get companies by owner id: entity not found",
		},
		{
			name:           "with other error",
			ownerId:        1,
			want:           nil,
			mockEntities:   nil,
			mockError:      errors.New("test error"),
			wantErr:        true,
			wantErrMessage: "usecases.GetCompaniesByOwnerId: get companies by owner id: test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			storageMock := mocks.NewMockcompanyStorage(mockCtrl)
			storageMock.EXPECT().GetCompaniesByOwnerId(gomock.Any(), tt.ownerId).Return(tt.mockEntities, tt.mockError).Times(1)

			u := &companyUsecases{
				cs: storageMock,
			}

			got, err := u.GetCompaniesByOwnerId(context.Background(), tt.ownerId)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("companyUsecases.GetCompaniesByOwnerId() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.Equal(t, tt.wantErrMessage, err.Error())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("companyUsecases.GetCompaniesByOwnerId() = %v, want %v", got, tt.want)
			}
		})
	}
}
