package entities

type Chat struct {
	Id           int64  `db:"id"`
	CompanyId    int64  `db:"company_id"`
	TelegramId   int64  `db:"telegram_id"`
}
