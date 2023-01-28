package models

import (
	"github.com/xybor-x/xypriv"
	"gorm.io/gorm"
)

const mask = "******************"

func NewMaskedUser() User {
	return User{
		Email:       mask,
		Username:    mask,
		FirstName:   mask,
		LastName:    mask,
		PhoneNumber: mask,
		Address:     mask,
		Role:        mask,
	}
}

var ReadableUserCols = map[string]any{
	"id": nil, "email": nil, "username": nil, "first_name": nil,
	"last_name": nil, "phone_number": nil, "address": nil, "role": nil,
}

var EditableUserCols = []string{
	"username", "first_name", "last_name", "phone_number", "address", "role",
}

var Roles = [...]string{"admin", "mod", "member"}

// User represents for a user.
type User struct {
	gorm.Model

	Email       string `gorm:"unique;not null"`
	Username    string `gorm:"unique;not null;index,unique"`
	FirstName   string
	LastName    string
	PhoneNumber string
	Address     string
	Role        string

	UserCredential UserCredential `gorm:"foreignKey:Email;references:Email;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Client         []Client       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (u User) IsReadable(col string) bool {
	_, ok := ReadableUserCols[col]
	return ok
}

// Relation implements xypriv.Subject interface.
func (u User) Relation(ctx any, s xypriv.Subject) xypriv.Relation {
	if another, ok := s.(User); ok {
		if u.Username == another.Username && u.Username != "" {
			return "self"
		}
		if u.ID == another.ID && u.ID != 0 {
			return "self"
		}
	}

	switch u.Role {
	case "admin":
		return "admin"
	case "mod":
		return "mod"
	}

	switch ctx {
	case nil:
		return "member"
	}

	return "anyone"
}

func (u User) Owner() xypriv.Subject {
	return u
}

func (u User) Context() any {
	return nil
}

func (u User) Permission(a ...string) xypriv.AccessLevel {
	if len(a) < 1 {
		return xypriv.NotSupport
	}

	switch a[0] {
	case "read":
		if len(a) != 2 {
			return xypriv.NotSupport
		}
		switch a[1] {
		case "username":
			return xypriv.LowPrivate
		case "id", "email":
			return xypriv.MediumPrivate
		case "phone_number", "address", "first_name", "last_name":
			return xypriv.LowConfidential
		}
	case "update":
		if len(a) != 2 {
			return xypriv.NotSupport
		}
		switch a[1] {
		case "first_name", "last_name", "username":
			return xypriv.TopSecret
		case "phone_number", "address":
			return xypriv.LowConfidential
		}
	}

	return xypriv.NotSupport
}
