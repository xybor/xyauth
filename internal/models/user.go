package models

import (
	"github.com/xybor-x/xypriv"
	"gorm.io/gorm"
)

// User represents for a user.
type User struct {
	gorm.Model `json:"-"`

	Email string `gorm:"unique;index:,unique"`
	Name  string
	Role  string

	Client []Client `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Relation implements xypriv.Subject interface.
func (u User) Relation(ctx any, s xypriv.Subject) xypriv.Relation {
	switch u.Role {
	case "admin":
		return "admin"
	case "mod":
		return "mod"
	}

	switch ctx {
	case nil:
		switch t := s.(type) {
		case User:
			if t.Email == u.Email {
				return "self"
			}
		}
		return "member"
	}

	return "anyone"
}
