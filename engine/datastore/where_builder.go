package datastore

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"gorm.io/gorm"
)

// whereBuilder holds a state during custom where preparation
type whereBuilder struct {
	client ClientInterface
	tx     *gorm.DB
	varNum int
}

// processConditions will process all conditions
func (builder *whereBuilder) processConditions(tx customWhereInterface, conditions map[string]interface{}, parentKey *string) {
	for key, condition := range conditions {
		switch {
		case key == conditionAnd:
			builder.processWhereAnd(tx, condition)
		case key == conditionOr:
			builder.processWhereOr(tx, condition)
		case key == conditionGreaterThan:
			builder.applyCondition(tx, *parentKey, ">", condition)
		case key == conditionLessThan:
			builder.applyCondition(tx, *parentKey, "<", condition)
		case key == conditionGreaterThanOrEqual:
			builder.applyCondition(tx, *parentKey, ">=", condition)
		case key == conditionLessThanOrEqual:
			builder.applyCondition(tx, *parentKey, "<=", condition)
		case key == conditionExists:
			builder.applyExistsCondition(tx, *parentKey, condition.(bool))
		case StringInSlice(key, builder.client.GetArrayFields()):
			builder.applyJSONArrayContains(tx, key, condition.(string))
		case StringInSlice(key, builder.client.GetObjectFields()):
			builder.applyJSONCondition(tx, key, condition)
		case condition == nil:
			builder.applyCondition(tx, key, "IS NULL", nil)
		default:
			v := reflect.ValueOf(condition)
			if v.Kind() == reflect.Map {
				dict := convertToDict(condition)
				builder.processConditions(tx, dict, &key)
			} else {
				builder.applyCondition(tx, key, "=", condition)
			}
		}
	}
}

func (builder *whereBuilder) applyCondition(tx customWhereInterface, key string, operator string, condition interface{}) {
	columnName := builder.getColumnNameOrPanic(key)

	if condition == nil {
		tx.Where(columnName + " " + operator)
		return
	}
	varName := builder.nextVarName()
	query := fmt.Sprintf("%s %s @%s", columnName, operator, varName)
	tx.Where(query, map[string]interface{}{varName: builder.formatCondition(condition)})
}

func (builder *whereBuilder) applyExistsCondition(tx customWhereInterface, key string, condition bool) {
	operator := "IS NULL"
	if condition {
		operator = "IS NOT NULL"
	}
	builder.applyCondition(tx, key, operator, nil)
}

// applyJSONArrayContains will apply array condition on JSON Array field - client.GetArrayFields()
func (builder *whereBuilder) applyJSONArrayContains(tx customWhereInterface, key string, condition string) {
	columnName := builder.getColumnNameOrPanic(key)

	engine := builder.client.Engine()

	switch engine {
	case PostgreSQL:
		builder.applyPostgresJSONB(tx, columnName, fmt.Sprintf(`["%s"]`, condition))
	case SQLite:
		varName := builder.nextVarName()
		tx.Where(
			fmt.Sprintf("EXISTS (SELECT 1 FROM json_each(%s) WHERE value = @%s)", columnName, varName),
			map[string]interface{}{varName: condition},
		)
	case Empty:
		panic("Database engine not configured")
	default:
		panic("Unknown database engine")
	}
}

// applyJSONCondition will apply condition on JSON Object field - client.GetObjectFields()
func (builder *whereBuilder) applyJSONCondition(tx customWhereInterface, key string, condition interface{}) {
	if isEmptyCondition(condition) {
		return
	}

	columnName := builder.getColumnNameOrPanic(key)
	engine := builder.client.Engine()

	if engine == PostgreSQL {
		builder.applyPostgresJSONB(tx, columnName, condition)
	} else if engine == SQLite {
		builder.applyJSONExtract(tx, columnName, condition)
	} else {
		panic("Database engine not supported")
	}
}

func (builder *whereBuilder) applyPostgresJSONB(tx customWhereInterface, columnName string, condition interface{}) {
	varName := builder.nextVarName()
	query := fmt.Sprintf("%s::jsonb @> @%s", columnName, varName)
	tx.Where(query, map[string]interface{}{varName: condition})
}

func (builder *whereBuilder) applyJSONExtract(tx customWhereInterface, columnName string, condition interface{}) {
	dict := convertToDict(condition)
	for key, value := range dict {
		keyVarName := builder.nextVarName()
		valueVarName := builder.nextVarName()
		query := fmt.Sprintf("JSON_EXTRACT(%s, @%s) = @%s", columnName, keyVarName, valueVarName)
		tx.Where(query, map[string]interface{}{
			keyVarName:   fmt.Sprintf("$.%s", key),
			valueVarName: value,
		})
	}
}

func (builder *whereBuilder) processWhereAnd(tx customWhereInterface, condition interface{}) {
	accumulator := &txAccumulator{
		WhereClauses: make([]string, 0),
		Vars:         make(map[string]interface{}),
	}
	for _, c := range condition.([]map[string]interface{}) {
		builder.processConditions(accumulator, c, nil)
	}

	query := " ( " + strings.Join(accumulator.WhereClauses, " AND ") + " ) "
	if len(accumulator.Vars) > 0 {
		tx.Where(query, accumulator.Vars)
	} else {
		tx.Where(query)
	}
}

func (builder *whereBuilder) processWhereOr(tx customWhereInterface, condition interface{}) {
	or := make([]string, 0)
	orVars := make(map[string]interface{})
	for _, cond := range condition.([]map[string]interface{}) {
		statement := make([]string, 0)
		accumulator := &txAccumulator{
			WhereClauses: make([]string, 0),
			Vars:         make(map[string]interface{}),
		}
		builder.processConditions(accumulator, cond, nil)
		statement = append(statement, accumulator.WhereClauses...)
		for varName, varValue := range accumulator.Vars {
			orVars[varName] = varValue
		}
		or = append(or, strings.Join(statement[:], " AND "))
	}

	query := " ( (" + strings.Join(or, ") OR (") + ") ) "
	if len(orVars) > 0 {
		tx.Where(query, orVars)
	} else {
		tx.Where(query)
	}
}

func (builder *whereBuilder) nextVarName() string {
	varName := fmt.Sprintf("var%d", builder.varNum)
	builder.varNum++
	return varName
}

func (builder *whereBuilder) getColumnNameOrPanic(key string) string {
	columnName, ok := GetColumnName(key, builder.tx.Statement.Model, builder.tx)
	if !ok {
		panic(fmt.Errorf("column %s does not exist in the model", key))
	}

	return columnName
}

func (builder *whereBuilder) formatCondition(condition interface{}) interface{} {
	if nullType, ok := condition.(customtypes.NullTime); ok {
		return builder.formatNullTime(nullType)
	}

	return condition
}

func (builder *whereBuilder) formatNullTime(condition customtypes.NullTime) interface{} {
	if !condition.Valid {
		return nil
	}
	engine := builder.client.Engine()
	if engine == PostgreSQL {
		return condition.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	return condition.Time.Format("2006-01-02T15:04:05.000Z")
}

func convertToDict(object interface{}) map[string]interface{} {
	if converted, ok := object.(map[string]interface{}); ok {
		return converted
	}
	vJSON, _ := json.Marshal(object)

	var converted map[string]interface{}
	_ = json.Unmarshal(vJSON, &converted)
	return converted
}

func isEmptyCondition(condition interface{}) bool {
	val := reflect.ValueOf(condition)
	for ; val.Kind() == reflect.Ptr; val = val.Elem() {
		if val.IsNil() {
			return true
		}
	}
	kind := val.Kind()
	if kind == reflect.Map || kind == reflect.Slice {
		return val.Len() == 0
	}

	return false
}
