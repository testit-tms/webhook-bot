package entities

// Owner represents the owner of a company.
type Owner struct {
	ID           int64  `db:"id"`
	TelegramID   int64  `db:"telegram_id"`
	TelegramName string `db:"telegram_name"`
}

// OwnerInfo represents information about the owner of a company.
type OwnerInfo struct {
	TelegramID   int64
	TelegramName string
}
