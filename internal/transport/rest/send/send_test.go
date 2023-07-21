package send

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/lib/handlers"
	"github.com/testit-tms/webhook-bot/internal/lib/logger/slogdiscard"
	"github.com/testit-tms/webhook-bot/internal/transport/rest/send/mocks"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		parseMode string
		chatIds   string
		respCode  int
		respError string
		mockTimes int
		mockError error
	}{
		{
			name:      "Success",
			message:   "test message",
			parseMode: "MarkdownV2",
			chatIds:   "12345",
			respCode:  http.StatusOK,
			respError: "message sent",
			mockTimes: 1,
		},
		{
			name:      "invalid request",
			message:   "test message",
			parseMode: "MarkdownV2",
			chatIds:   "\"12345\"",
			respCode:  http.StatusBadRequest,
			respError: "failed to decode request",
			mockTimes: 0,
		},
		{
			name:      "invalid parse mode",
			message:   "test message",
			parseMode: "qwerty",
			chatIds:   "12345",
			respCode:  http.StatusBadRequest,
			respError: "field ParseMode must be empty or have following value: markdownv2, markdown or html",
			mockTimes: 0,
		},
		{
			name:      "empty parse mode",
			message:   "test message",
			chatIds:   "12345",
			respCode:  http.StatusOK,
			respError: "message sent",
			mockTimes: 1,
		},
		{
			name:      "empty message",
			chatIds:   "12345",
			respCode:  http.StatusBadRequest,
			respError: "field Message is a required field",
			mockTimes: 0,
		},
		{
			name:      "empty chatids",
			message:   "test message",
			respCode:  http.StatusOK,
			respError: "message sent",
			mockTimes: 1,
		},
		{
			name:      "error send message",
			message:   "test message",
			respCode:  http.StatusInternalServerError,
			respError: "can't send message",
			mockTimes: 1,
			mockError: errors.New("some error"),
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			senderMock := mocks.NewMocksender(mockCtrl)

			mes := entities.Message{
				Text:      tc.message,
				ParseMode: entities.ParseString(tc.parseMode),
			}

			if tc.mockError != nil {
				senderMock.EXPECT().SendMessage(gomock.Any(), mes).Return(tc.mockError).Times(tc.mockTimes)
			} else {
				senderMock.EXPECT().SendMessage(gomock.Any(), mes).Return(nil).Times(tc.mockTimes)
			}

			handler := New(slogdiscard.NewDiscardLogger(), senderMock)

			input := fmt.Sprintf(`{"message": "%s", "parseMode": "%s", "chatIds": ["%s"]}`, tc.message, tc.parseMode, tc.chatIds)

			req, err := http.NewRequest(http.MethodPost, "/telegram", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.respCode)

			body := rr.Body.String()

			var respMessage string

			if tc.respCode == http.StatusOK {
				respMessage = "message sent"
			} else {
				var resp handlers.ErrorResponse
				require.NoError(t, json.Unmarshal([]byte(body), &resp))

				respMessage = resp.Message
			}

			require.Equal(t, tc.respError, respMessage)
		})
	}
}
