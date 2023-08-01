package send

import "github.com/testit-tms/webhook-bot/internal/entities"

type Request struct {
	Message   string  `json:"message" validate:"required"`
	ParseMode string  `json:"parseMode,omitempty" validate:"parse-mode"`
	ChatIds   []int64 `json:"chatIds,omitempty"`
}

func (m *Request) convertToDomain() entities.Message {
	return entities.Message{
		ParseMode: entities.ParseString(m.ParseMode),
		Text:      m.Message,
		ChatIds:   m.ChatIds,
	}
}
