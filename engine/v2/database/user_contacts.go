package database

import (
	"gorm.io/gorm"
)

// UserContact represents a contact between two users but is assigned to one.
type UserContact struct {
	gorm.Model

	FullName string
	Status   string
	Paymail  string
	PubKey   string

	UserID string
	User   *User `gorm:"foreignKey:UserID"`
}
