package database

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Address represents a user's (bitcoin) addresses.
type Address struct {
	Address string `gorm:"type:char(34);primaryKey"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CustomInstructions datatypes.JSONSlice[CustomInstruction]

	UserID string
	User   *User `gorm:"foreignKey:UserID"`
}

// CustomInstruction represents a custom instruction how to unlock a UTXO.
type CustomInstruction struct {
	Type        string
	Instruction string
}
