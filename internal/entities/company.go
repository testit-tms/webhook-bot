package entities

type Company struct {
	Id      int    `db:"id"`
	OwnerId int    `db:"owner_id"`
	Token   string `db:"token"`
	Name    string `db:"name"`
	Email   string `db:"email"`
}

type CompanyInfo struct {
	Name  string
	Email string
	Owner OwnerInfo
}
