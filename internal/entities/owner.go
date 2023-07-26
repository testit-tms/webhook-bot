package entities

type Owner struct {
	Id           int64  `db:"id"`
	TelegramId   int64  `db:"telegram_id"`
	TelegramName string `db:"telegram_name"`
}

type OwnerInfo struct {
	TelegramId   int64
	TelegramName string
}
