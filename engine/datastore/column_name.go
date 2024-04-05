package datastore

import (
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var modelsCache = sync.Map{}

// GetColumnName checks if the model has provided columnName as DBName or (struct field)Name
// Returns (DBName, true) if the column exists otherwise (_, false)
// Uses global cache store (thread safe)
// Checking is case-sensitive
// The gdb param is optional. When is provided, the actual naming strategy is used; otherwise default
func GetColumnName(columnName string, model interface{}, gdb *gorm.DB) (string, bool) {
	var namer schema.Namer
	if gdb != nil {
		namer = gdb.NamingStrategy
	} else {
		namer = schema.NamingStrategy{}
	}

	sch, err := schema.Parse(model, &modelsCache, namer)
	if err != nil {
		return "", false
	}
	if field, ok := sch.FieldsByDBName[columnName]; ok {
		return field.DBName, true
	}

	if field, ok := sch.FieldsByName[columnName]; ok {
		return field.DBName, true
	}

	return "", false
}
