package repository

import "gorm.io/gorm"

func mapConditionsToScopes(conditions map[string]interface{}) []func(tx *gorm.DB) *gorm.DB {
	scopes := make([]func(tx *gorm.DB) *gorm.DB, 0, len(conditions))
	for key, value := range conditions {
		scopes = append(scopes, func(key string, value interface{}) func(tx *gorm.DB) *gorm.DB {
			return func(tx *gorm.DB) *gorm.DB {
				return tx.Where(key+" = ?", value)
			}
		}(key, value))
	}

	return scopes
}
