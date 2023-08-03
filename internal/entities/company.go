package entities

type Company struct {
	Id      int64  `db:"id"`
	OwnerId int64  `db:"owner_id"`
	Token   string `db:"token"`
	Name    string `db:"name"`
	Email   string `db:"email"`
}

type CompanyRegistrationInfo struct {
	Name  string
	Email string
	Owner OwnerInfo
}

type CompanyInfo struct {
	Id      int64
	OwnerId int64
	Token   string
	Name    string
	Email   string
	ChatIds []int64
}
