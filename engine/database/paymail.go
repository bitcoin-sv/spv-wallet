package database

import "gorm.io/gorm"

// Paymail represents a paymail address
type Paymail struct {
	gorm.Model

	Alias      string
	Domain     string
	PublicName string
	AvatarURL  string

	UserID string
	User   *User `gorm:"foreignKey:UserID"`
}
