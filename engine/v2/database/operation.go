package database

import "time"

// Operation represents a user's operation on a transaction.
type Operation struct {
	TxID   string `gorm:"primaryKey"`
	UserID string `gorm:"primaryKey"`

	CreatedAt time.Time

	Counterparty string
	Type         string
	Value        int64

	User        *User               `gorm:"foreignKey:UserID"`
	Transaction *TrackedTransaction `gorm:"foreignKey:TxID"`
}
