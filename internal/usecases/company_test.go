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
		name             string
		ownerId          int64
		want             entities.CompanyInfo
		mockCompEntities entities.Company
		mockCompError    error
		mockCompTimes    int
		mockChatEntities []entities.Chat
		mockChatError    error
		mockChatTimes    int
		wantErr          bool
		wantErrMessage   string
	}{
		{
			name:    "success",
			ownerId: 1,
			want: entities.CompanyInfo{
				ID:      12,
				OwnerID: 21,
				Token:   "token",
				Name:    "Yandex",
				Email:   "info@ya.ru",
				ChatIds: []int64{
					123,
				},
			},
			mockCompEntities: entities.Company{
				ID:      12,
				OwnerID: 21,
				Token:   "token",
				Name:    "Yandex",
				Email:   "info@ya.ru",
			},
			mockCompError: nil,
			mockCompTimes: 1,
			mockChatEntities: []entities.Chat{
				{
					Id:         1,
					CompanyID:  12,
					TelegramID: 123,
				},
			},
			mockChatError: nil,
			mockChatTimes: 1,
			wantErr:       false,
		},
		{
			name:             "company with ErrNotFound",
			ownerId:          1,
			want:             entities.CompanyInfo{},
			mockCompEntities: entities.Company{},
			mockCompError:    storage.ErrNotFound,
			mockCompTimes:    1,
			mockChatTimes:    0,
			wantErr:          true,
			wantErrMessage:   "usecases.GetCompanyByOwnerTelegramId: company not found",
		},
		{
			name:             "company with other error",
			ownerId:          1,
			want:             entities.CompanyInfo{},
			mockCompEntities: entities.Company{},
			mockCompError:    errors.New("test error"),
			mockCompTimes:    1,
			mockChatTimes:    0,
			wantErr:          true,
			wantErrMessage:   "usecases.GetCompanyByOwnerTelegramId: get company by owner id: test error",
		},
		{
			name:    "chats with ErrNotFound",
			ownerId: 1,
			want: entities.CompanyInfo{
				ID:      12,
				OwnerID: 21,
				Token:   "token",
				Name:    "Yandex",
				Email:   "info@ya.ru",
				ChatIds: nil,
			},
			mockCompEntities: entities.Company{
				ID:      12,
				OwnerID: 21,
				Token:   "token",
				Name:    "Yandex",
				Email:   "info@ya.ru",
			},
			mockCompError:    nil,
			mockCompTimes:    1,
			mockChatEntities: []entities.Chat{},
			mockChatError:    storage.ErrNotFound,
			mockChatTimes:    1,
			wantErr:          false,
		},
		{
			name:    "chats with other error",
			ownerId: 1,
			want: entities.CompanyInfo{
				ID:      12,
				OwnerID: 21,
				Token:   "token",
				Name:    "Yandex",
				Email:   "info@ya.ru",
				ChatIds: nil,
			},
			mockCompEntities: entities.Company{
				ID:      12,
				OwnerID: 21,
				Token:   "token",
				Name:    "Yandex",
				Email:   "info@ya.ru",
			},
			mockCompError:    nil,
			mockCompTimes:    1,
			mockChatEntities: []entities.Chat{},
			mockChatError:    errors.New("test error"),
			mockChatTimes:    1,
			wantErr:          true,
			wantErrMessage:   "usecases.GetCompanyByOwnerTelegramId: get chats by company id: test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			companyMock := mocks.NewMockcompanyStorage(mockCtrl)
			companyMock.EXPECT().GetCompanyByOwnerTelegramId(gomock.Any(), tt.ownerId).Return(tt.mockCompEntities, tt.mockCompError).Times(tt.mockCompTimes)

			chatMock := mocks.NewMockchatStorage(mockCtrl)
			if tt.mockChatTimes != 0 {
				chatMock.EXPECT().GetChatsByCompanyId(gomock.Any(), tt.mockCompEntities.ID).Return(tt.mockChatEntities, tt.mockChatError).Times(tt.mockChatTimes)
			}

			u := NewCompanyUsecases(companyMock, chatMock)

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

func Test_companyUsecases_UpdateToken(t *testing.T) {
	tests := []struct {
		name                 string
		ownerId              int64
		wantErr              bool
		mockCompEntities     entities.Company
		mockCompError        error
		mockCompTimes        int
		mockUpdateTokenError error
		mockUpdateTokenTimes int
	}{
		{
			name:    "success",
			ownerId: 1,
			wantErr: false,
			mockCompEntities: entities.Company{
				ID:      12,
				OwnerID: 21,
				Token:   "token",
				Name:    "Yandex",
				Email:   "",
			},
			mockCompError:        nil,
			mockCompTimes:        1,
			mockUpdateTokenError: nil,
			mockUpdateTokenTimes: 1,
		},
		{
			name:             "company with ErrNotFound",
			ownerId:          1,
			wantErr:          true,
			mockCompEntities: entities.Company{},
			mockCompError:    storage.ErrNotFound,
			mockCompTimes:    1,
		},
		{
			name:             "company with other error",
			ownerId:          1,
			wantErr:          true,
			mockCompEntities: entities.Company{},
			mockCompError:    errors.New("test error"),
			mockCompTimes:    1,
		},
		{
			name:    "update token with error",
			wantErr: true,
			mockCompEntities: entities.Company{
				ID:      12,
				OwnerID: 21,
				Token:   "token",
				Name:    "Yandex",
				Email:   "",
			},
			mockCompError:        nil,
			mockCompTimes:        1,
			mockUpdateTokenError: errors.New("test error"),
			mockUpdateTokenTimes: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			companyMock := mocks.NewMockcompanyStorage(mockCtrl)
			companyMock.EXPECT().GetCompanyByOwnerTelegramId(gomock.Any(), tt.ownerId).Return(tt.mockCompEntities, tt.mockCompError).Times(tt.mockCompTimes)
			if tt.mockCompError == nil {
				companyMock.EXPECT().UpdateToken(gomock.Any(), tt.mockCompEntities.ID, gomock.Any()).Return(tt.mockUpdateTokenError).Times(tt.mockUpdateTokenTimes)
			}

			u := NewCompanyUsecases(companyMock, nil)

			err := u.UpdateToken(context.Background(), tt.ownerId)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("companyUsecases.UpdateToken() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}
