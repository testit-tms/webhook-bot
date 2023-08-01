package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/lib/logger/slogdiscard"
	"github.com/testit-tms/webhook-bot/internal/storage"
	"github.com/testit-tms/webhook-bot/internal/usecases/mocks"
	"go.uber.org/mock/gomock"
)

func Test_sendMessageUsacases_SendMessage(t *testing.T) {
	tests := []struct {
		name             string
		msg              entities.Message
		mockChatEntities []entities.Chat
		mockChatError    error
		mockChatTimes    int
		mockBotEntity    entities.Message
		mockBotError     error
		mockBotTimes     int
		wantErr          bool
		wantErrMessage   string
	}{
		{
			name: "success",
			msg: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				Token:     "token",
				ChatIds:   []int64{123},
			},
			mockChatEntities: []entities.Chat{
				{
					Id:         1,
					TelegramId: 123,
					CompanyId:  12,
				},
			},
			mockChatError: nil,
			mockChatTimes: 1,
			mockBotEntity: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				ChatIds:   []int64{123},
				Token:     "token",
			},
			mockBotError: nil,
			mockBotTimes: 1,
			wantErr:      false,
		},
		{
			name: "get chats error sql not found",
			msg: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				Token:     "token",
				ChatIds:   []int64{123},
			},
			mockChatEntities: []entities.Chat{},
			mockChatError:    storage.ErrNotFound,
			mockChatTimes:    1,
			mockBotEntity: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				ChatIds:   []int64{123},
				Token:     "token",
			},
			mockBotError:   nil,
			mockBotTimes:   0,
			wantErr:        true,
			wantErrMessage: "usecases.SendMessage: chats not found: chats not found",
		},
		{
			name: "get chats other error",
			msg: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				Token:     "token",
				ChatIds:   []int64{123},
			},
			mockChatEntities: []entities.Chat{},
			mockChatError:    errors.New("error"),
			mockChatTimes:    1,
			mockBotEntity: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				ChatIds:   []int64{123},
				Token:     "token",
			},
			mockBotError:   nil,
			mockBotTimes:   0,
			wantErr:        true,
			wantErrMessage: "usecases.SendMessage: get chats by company token: chats not found",
		},
		{
			name: "success without chat ids",
			msg: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				Token:     "token",
				ChatIds:   []int64{},
			},
			mockChatEntities: []entities.Chat{
				{
					Id:         1,
					TelegramId: 123,
					CompanyId:  12,
				},
			},
			mockChatError: nil,
			mockChatTimes: 1,
			mockBotEntity: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				ChatIds:   []int64{123},
				Token:     "token",
			},
			mockBotError: nil,
			mockBotTimes: 1,
			wantErr:      false,
		},
		{
			name: "without chat ids error",
			msg: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				Token:     "token",
				ChatIds:   []int64{},
			},
			mockChatEntities: []entities.Chat{
				{
					Id:         1,
					TelegramId: 123,
					CompanyId:  12,
				},
			},
			mockChatError: nil,
			mockChatTimes: 1,
			mockBotEntity: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				ChatIds:   []int64{123},
				Token:     "token",
			},
			mockBotError:   errors.New("error"),
			mockBotTimes:   1,
			wantErr:        true,
			wantErrMessage: "usecases.SendMessage: can not send message: can not send message",
		},
		{
			name: "not allowed chats error",
			msg: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				Token:     "token",
				ChatIds:   []int64{321},
			},
			mockChatEntities: []entities.Chat{
				{
					Id:         1,
					TelegramId: 123,
					CompanyId:  12,
				},
			},
			mockChatError: nil,
			mockChatTimes: 1,
			mockBotEntity: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				ChatIds:   []int64{123},
				Token:     "token",
			},
			mockBotError:   errors.New("error"),
			mockBotTimes:   0,
			wantErr:        true,
			wantErrMessage: "usecases.SendMessage: chats not allowed: chats not allowed",
		},
		{
			name: "with chat ids error",
			msg: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				Token:     "token",
				ChatIds:   []int64{123},
			},
			mockChatEntities: []entities.Chat{
				{
					Id:         1,
					TelegramId: 123,
					CompanyId:  12,
				},
			},
			mockChatError: nil,
			mockChatTimes: 1,
			mockBotEntity: entities.Message{
				Text:      "text",
				ParseMode: entities.MarkdownV2,
				ChatIds:   []int64{123},
				Token:     "token",
			},
			mockBotError:   errors.New("error"),
			mockBotTimes:   1,
			wantErr:        true,
			wantErrMessage: "usecases.SendMessage: can not send message: can not send message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockChat := mocks.NewMockchatGeter(ctrl)
			mockChat.EXPECT().GetChatsByCompanyToken(gomock.Any(), tt.msg.Token).Return(tt.mockChatEntities, tt.mockChatError).Times(tt.mockChatTimes)

			mockBot := mocks.NewMockbotSender(ctrl)
			if tt.mockBotTimes != 0 {
				mockBot.EXPECT().SendMessage(gomock.Any(), tt.mockBotEntity).Return(tt.mockBotError).Times(tt.mockBotTimes)
			}

			u := NewSendMessageUsecases(slogdiscard.NewDiscardLogger(), mockChat, mockBot)

			if err := u.SendMessage(context.Background(), tt.msg); err != nil {
				if !tt.wantErr {
					t.Errorf("sendMessageUsacases.SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				}

				assert.Equal(t, tt.wantErrMessage, err.Error())
			}
		})
	}
}
