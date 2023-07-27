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
