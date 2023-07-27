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

func Test_companyUsecases_GetCompanyByOwnerId(t *testing.T) {
	tests := []struct {
		name           string
		ownerId        int64
		want           entities.Company
		mockEntities   entities.Company
		mockError      error
		wantErr        bool
		wantErrMessage string
	}{
		{
			name:    "success",
			ownerId: 1,
			want: entities.Company{
				Id:      12,
				OwnerId: 21,
				Token:   "token",
				Name:    "Yandex",
				Email:   "info@ya.ru",
			},
			mockEntities: entities.Company{
				Id:      12,
				OwnerId: 21,
				Token:   "token",
				Name:    "Yandex",
				Email:   "info@ya.ru",
			},
			mockError: nil,
			wantErr:   false,
		},
		{
			name:           "with ErrNotFound",
			ownerId:        1,
			want:           entities.Company{},
			mockEntities:   entities.Company{},
			mockError:      storage.ErrNotFound,
			wantErr:        true,
			wantErrMessage: "usecases.GetCompanyByOwnerTelegramId: company not found",
		},
		{
			name:           "with other error",
			ownerId:        1,
			want:           entities.Company{},
			mockEntities:   entities.Company{},
			mockError:      errors.New("test error"),
			wantErr:        true,
			wantErrMessage: "usecases.GetCompanyByOwnerTelegramId: get company by owner id: test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			storageMock := mocks.NewMockcompanyStorage(mockCtrl)
			storageMock.EXPECT().GetCompanyByOwnerTelegramId(gomock.Any(), tt.ownerId).Return(tt.mockEntities, tt.mockError).Times(1)

			u := &companyUsecases{
				cs: storageMock,
			}

			got, err := u.GetCompanyByOwnerTelegramId(context.Background(), tt.ownerId)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("companyUsecases.GetCompanyByOwnerTelegramId() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.Equal(t, tt.wantErrMessage, err.Error())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("companyUsecases.GetCompanyByOwnerTelegramId() = %v, want %v", got, tt.want)
			}
		})
	}
}
