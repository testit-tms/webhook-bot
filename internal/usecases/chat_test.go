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

func Test_chatUsecases_AddChat(t *testing.T) {
	tests := []struct {
		name           string
		chat           entities.Chat
		want           entities.Chat
		wantErr        bool
		wantErrMessage string
		wantError      error
	}{
		{
			name: "success",
			chat: entities.Chat{
				CompanyId:  12,
				TelegramId: 123,
			},
			want: entities.Chat{
				Id:         1,
				CompanyId:  12,
				TelegramId: 123,
			},
			wantErr:        false,
			wantErrMessage: "",
			wantError:      nil,
		},
		{
			name: "error",
			chat: entities.Chat{
				CompanyId:  12,
				TelegramId: 123,
			},
			want:           entities.Chat{},
			wantErr:        true,
			wantErrMessage: "error",
			wantError:      errors.New("error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			chatMock := mocks.NewMockchatsStorage(mockCtrl)
			chatMock.EXPECT().AddChat(gomock.Any(), tt.chat).Return(tt.want, tt.wantError)

			companyMock := mocks.NewMockcompanyStorage(mockCtrl)
			u := NewChatUsecases(chatMock, companyMock)

			got, err := u.AddChat(context.Background(), tt.chat)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("chatUsecases.AddChat() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.Equal(t, tt.wantErrMessage, err.Error())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("chatUsecases.AddChat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_chatUsecases_DeleteChatByTelegramId(t *testing.T) {
	tests := []struct {
		name                string
		ownerId             int64
		chatId              int64
		mockCompEntities    entities.Company
		mockCompError       error
		mockCompTimes       int
		mockGetChatEntities []entities.Chat
		mockGetChatError    error
		mockGetChatTimes    int
		mockDeleteChatError error
		mockDeleteChatTimes int
		wantErr             bool
		wantErrMessage      string
		wantError           error
	}{
		{
			name:    "success",
			ownerId: 1,
			chatId:  123,
			mockCompEntities: entities.Company{
				Id: 12,
			},
			mockCompError: nil,
			mockCompTimes: 1,
			mockGetChatEntities: []entities.Chat{
				{
					Id:         1,
					CompanyId:  12,
					TelegramId: 123,
				},
			},
			mockGetChatError:    nil,
			mockGetChatTimes:    1,
			mockDeleteChatError: nil,
			mockDeleteChatTimes: 1,
			wantErr:             false,
		},
		{
			name:             "company not found",
			ownerId:          1,
			chatId:           123,
			mockCompEntities: entities.Company{},
			mockCompError:    storage.ErrNotFound,
			mockCompTimes:    1,
			mockGetChatTimes: 0,
			wantErr:          true,
			wantErrMessage:   "usecases.DeleteChatByTelegramId: get company by owner id: entity not found",
		},
		{
			name:    "chat not found",
			ownerId: 1,
			chatId:  123,
			mockCompEntities: entities.Company{
				Id: 12,
			},
			mockCompError: nil,
			mockCompTimes: 1,
			mockGetChatEntities: []entities.Chat{
				{
					Id:         1,
					CompanyId:  12,
					TelegramId: 321,
				},
			},
			mockGetChatError:    nil,
			mockGetChatTimes:    1,
			mockDeleteChatError: nil,
			mockDeleteChatTimes: 0,
			wantErr:             true,
			wantErrMessage:      "usecases.DeleteChatByTelegramId: chat not found",
		},
		{
			name:    "get chats error",
			ownerId: 1,
			chatId:  123,
			mockCompEntities: entities.Company{
				Id: 12,
			},
			mockCompError:       nil,
			mockCompTimes:       1,
			mockGetChatEntities: []entities.Chat{},
			mockGetChatError:    errors.New("error"),
			mockGetChatTimes:    1,
			mockDeleteChatError: nil,
			mockDeleteChatTimes: 0,
			wantErr:             true,
			wantErrMessage:      "usecases.DeleteChatByTelegramId: get chats by company id: error",
		},
		{
			name:    "delete chat error",	
			ownerId: 1,
			chatId:  123,
			mockCompEntities: entities.Company{
				Id: 12,
			},
			mockCompError: nil,
			mockCompTimes: 1,
			mockGetChatEntities: []entities.Chat{
				{
					Id:         1,
					CompanyId:  12,
					TelegramId: 123,
				},
			},
			mockGetChatError:    nil,
			mockGetChatTimes:    1,
			mockDeleteChatError: errors.New("error"),
			mockDeleteChatTimes: 1,
			wantErr:             true,
			wantErrMessage:      "usecases.DeleteChatByTelegramId: delete chat by id: error",
		},
		

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)

			companyMock := mocks.NewMockcompanyStorage(mockCtrl)
			companyMock.EXPECT().GetCompanyByOwnerTelegramId(gomock.Any(), tt.ownerId).Return(tt.mockCompEntities, tt.mockCompError).Times(tt.mockCompTimes)

			chatMock := mocks.NewMockchatsStorage(mockCtrl)
			if tt.mockGetChatTimes != 0 {
				chatMock.EXPECT().GetChatsByCompanyId(gomock.Any(), tt.mockCompEntities.Id).Return(tt.mockGetChatEntities, tt.mockGetChatError).Times(tt.mockGetChatTimes)
			}

			if tt.mockDeleteChatTimes != 0 {
				chatMock.EXPECT().DeleteChatById(gomock.Any(), tt.mockGetChatEntities[0].Id).Return(tt.mockDeleteChatError).Times(tt.mockDeleteChatTimes)
			}

			u := NewChatUsecases(chatMock, companyMock)

			err := u.DeleteChatByTelegramId(context.Background(), tt.ownerId, tt.chatId)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("companyUsecases.GetCompanyByOwnerTelegramId() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.Equal(t, tt.wantErrMessage, err.Error())
			}
		})
	}
}
