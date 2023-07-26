package entities

type Chat struct {
	Id           int    `db:"id"`
	CompanyId    int    `db:"company_id"`
	TelegramId   int64  `db:"telegram_id"`
	TelegramName string `db:"telegram_name"`
}
