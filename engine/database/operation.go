package database

import "time"

type Operation struct {
	TxID   string `gorm:"primaryKey"`
	UserID string `gorm:"primaryKey"`

	CreatedAt time.Time

	Type  string
	Value int64

	User        *User               `gorm:"foreignKey:UserID"`
	Transaction *TrackedTransaction `gorm:"foreignKey:TxID"`
}
