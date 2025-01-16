package dbquery

import "gorm.io/gorm"

// UserID is a scope function that filters by user ID.
func UserID(id string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", id)
	}
}
