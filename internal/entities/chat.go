package entities

// Chat represents a chat entity.
type Chat struct {
	Id         int64 `db:"id"`
	CompanyID  int64 `db:"company_id"`
	TelegramID int64 `db:"telegram_id"`
}
