package models

// UserCredential represents for credential info of a User.
type UserCredential struct {
	Email    string `gorm:"unique;index:,unique"`
	Password string
}
