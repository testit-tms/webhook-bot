package storage

type Chat struct {
	Id           int    `db:"id"`
	CompanyId    int    `db:"company_id"`
	TelegramId   string `db:"telegram_id"`
	TelegramName string `db:"telegram_name"`
}
