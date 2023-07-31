package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testit-tms/webhook-bot/internal/entities"
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

			u := NewChatUsecases(chatMock)

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
