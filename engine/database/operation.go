package database

type Operation struct {
	TxID   string `gorm:"primaryKey"`
	UserID string `gorm:"primaryKey"`

	Type  string
	Value int64

	User        *User               `gorm:"foreignKey:UserID"`
	Transaction *TrackedTransaction `gorm:"foreignKey:TxID"`
}
