package storage

type Owner struct {
	Id           int    `db:"id"`
	TelegramId   string `db:"telegram_id"`
	TelegramName string `db:"telegram_name"`
}
