package database

import "gorm.io/gorm"

// Paymail represents a paymail address
type Paymail struct {
	gorm.Model

	Alias  string `gorm:"uniqueIndex:idx_alias_domain"`
	Domain string `gorm:"uniqueIndex:idx_alias_domain"`

	PublicName string
	Avatar     string

	UserID string
	User   *User `gorm:"foreignKey:UserID"`
}
