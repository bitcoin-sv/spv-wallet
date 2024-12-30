package database

import "gorm.io/gorm"

type Paymail struct {
	gorm.Model

	Alias      string
	Domain     string
	PublicName string
	AvatarURL  string

	UserID string
	User   *User `gorm:"foreignKey:UserID"`
}
