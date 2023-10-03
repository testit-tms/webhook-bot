package entities

// Company represents a company entity.
type Company struct {
	ID      int64  `db:"id"`
	OwnerID int64  `db:"owner_id"`
	Token   string `db:"token"`
	Name    string `db:"name"`
	Email   string `db:"email"`
}

// CompanyRegistrationInfo represents the information needed to register a new company.
type CompanyRegistrationInfo struct {
	Name  string
	Email string
	Owner OwnerInfo
}

// CompanyInfo represents the information about a company.
type CompanyInfo struct {
	ID      int64
	OwnerID int64
	Token   string
	Name    string
	Email   string
	ChatIds []int64
}
