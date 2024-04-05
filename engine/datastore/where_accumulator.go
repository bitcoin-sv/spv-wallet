package datastore

import (
	"gorm.io/gorm"
)

// customWhereInterface with single method Where which aligns with gorm.DB.Where
type customWhereInterface interface {
	Where(query interface{}, args ...interface{}) *gorm.DB
}

// txAccumulator holds the state of the nested conditions for recursive processing
type txAccumulator struct {
	WhereClauses []string
	Vars         map[string]interface{}
}

// Where makes txAccumulator implement customWhereInterface which will overload gorm.DB.Where behavior
func (tx *txAccumulator) Where(query interface{}, args ...interface{}) *gorm.DB {
	tx.WhereClauses = append(tx.WhereClauses, query.(string))

	if len(args) > 0 {
		for _, variables := range args {
			for key, value := range variables.(map[string]interface{}) {
				tx.Vars[key] = value
			}
		}
	}

	return nil
}
