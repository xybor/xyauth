package securitytoken

type AccessToken struct {
	UserID    uint
	Email     string
	Username  string
	FirstName string
	LastName  string
	Role      string
}

type RefreshToken struct {
	UserID   uint
	Family   string
	FamilyID uint
}
