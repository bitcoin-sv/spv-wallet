package datastore

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	customtypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
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

// whereBuilder holds a state during custom where preparation
type whereBuilder struct {
	client ClientInterface
	tx     *gorm.DB
	varNum int
}

// ApplyCustomWhere adds conditions to the gorm db instance
// it returns a tx of type *gorm.DB with a model and conditions applied
func ApplyCustomWhere(client ClientInterface, gdb *gorm.DB, conditions map[string]interface{}, model interface{}) (tx *gorm.DB, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error processing conditions: %v", r)
		}
	}()

	tx = gdb.Model(model)

	builder := &whereBuilder{
		client: client,
		tx:     tx,
		varNum: 0,
	}

	builder.processConditions(tx, conditions, nil)
	return
}

func (builder *whereBuilder) nextVarName() string {
	varName := "var" + strconv.Itoa(builder.varNum)
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

func (builder *whereBuilder) applyCondition(tx customWhereInterface, key string, operator string, condition interface{}) {
	columnName := builder.getColumnNameOrPanic(key)

	varName := builder.nextVarName()
	query := fmt.Sprintf("%s %s @%s", columnName, operator, varName)
	tx.Where(query, map[string]interface{}{varName: builder.formatCondition(condition)})
}

func (builder *whereBuilder) applyExistsCondition(tx customWhereInterface, key string, condition bool) {
	columnName := builder.getColumnNameOrPanic(key)

	operator := "IS NULL"
	if condition {
		operator = "IS NOT NULL"
	}
	tx.Where(columnName + " " + operator)
}

// processConditions will process all conditions
func (builder *whereBuilder) processConditions(tx customWhereInterface, conditions map[string]interface{}, parentKey *string) {
	for key, condition := range conditions {
		if key == conditionAnd {
			builder.processWhereAnd(tx, condition)
		} else if key == conditionOr {
			builder.processWhereOr(tx, conditions[conditionOr])
		} else if key == conditionGreaterThan {
			builder.applyCondition(tx, *parentKey, ">", condition)
		} else if key == conditionLessThan {
			builder.applyCondition(tx, *parentKey, "<", condition)
		} else if key == conditionGreaterThanOrEqual {
			builder.applyCondition(tx, *parentKey, ">=", condition)
		} else if key == conditionLessThanOrEqual {
			builder.applyCondition(tx, *parentKey, "<=", condition)
		} else if key == conditionExists {
			builder.applyExistsCondition(tx, *parentKey, condition.(bool))
		} else if StringInSlice(key, builder.client.GetArrayFields()) {
			tx.Where(builder.whereSlice(key, builder.formatCondition(condition)))
		} else if StringInSlice(key, builder.client.GetObjectFields()) {
			tx.Where(builder.whereObject(key, builder.formatCondition(condition)))
		} else {
			if condition == nil {
				tx.Where(key + " IS NULL")
			} else {
				v := reflect.ValueOf(condition)
				switch v.Kind() { //nolint:exhaustive // not all cases are needed
				case reflect.Map:
					if _, ok := condition.(map[string]interface{}); ok {
						builder.processConditions(tx, condition.(map[string]interface{}), &key) //nolint:scopelint // ignore for now
					} else {
						c, _ := json.Marshal(condition) //nolint:errchkjson // this check might break the current code
						var cc map[string]interface{}
						_ = json.Unmarshal(c, &cc)
						builder.processConditions(tx, cc, &key) //nolint:scopelint // ignore for now
					}
				default:
					builder.applyCondition(tx, key, "=", condition)
				}
			}
		}
	}
}

// formatCondition will format the conditions
func (builder *whereBuilder) formatCondition(condition interface{}) interface{} {
	switch v := condition.(type) {
	case customtypes.NullTime:
		if v.Valid {
			engine := builder.client.Engine()
			if engine == MySQL {
				return v.Time.Format("2006-01-02 15:04:05")
			} else if engine == PostgreSQL {
				return v.Time.Format("2006-01-02T15:04:05Z07:00")
			}
			// default & SQLite
			return v.Time.Format("2006-01-02T15:04:05.000Z")
		}
		return nil
	}

	return condition
}

// processWhereAnd will process the AND statements
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

// processWhereOr will process the OR statements
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

// escapeDBString will escape the database string
func escapeDBString(s string) string {
	rs := strings.Replace(s, "'", "\\'", -1)
	return strings.Replace(rs, "\"", "\\\"", -1)
}

// whereObject generates the where object
func (builder *whereBuilder) whereObject(k string, v interface{}) string {
	queryParts := make([]string, 0)

	// we don't know the type, we handle the rangeValue as a map[string]interface{}
	vJSON, _ := json.Marshal(v) //nolint:errchkjson // this check might break the current code

	var rangeV map[string]interface{}
	_ = json.Unmarshal(vJSON, &rangeV)

	engine := builder.client.Engine()

	for rangeKey, rangeValue := range rangeV {
		if engine == MySQL || engine == SQLite {
			switch vv := rangeValue.(type) {
			case string:
				rangeValue = "\"" + escapeDBString(rangeValue.(string)) + "\""
				queryParts = append(queryParts, "JSON_EXTRACT("+k+", '$."+rangeKey+"') = "+rangeValue.(string))
			default:
				metadataJSON, _ := json.Marshal(vv) //nolint:errchkjson // this check might break the current code
				var metadata map[string]interface{}
				_ = json.Unmarshal(metadataJSON, &metadata)
				for kk, vvv := range metadata {
					mJSON, _ := json.Marshal(vvv) //nolint:errchkjson // this check might break the current code
					vvv = string(mJSON)
					queryParts = append(queryParts, "JSON_EXTRACT("+k+", '$."+rangeKey+"."+kk+"') = "+vvv.(string))
				}
			}
		} else if engine == PostgreSQL {
			switch vv := rangeValue.(type) {
			case string:
				rangeValue = "\"" + escapeDBString(rangeValue.(string)) + "\""
			default:
				metadataJSON, _ := json.Marshal(vv) //nolint:errchkjson // this check might break the current code
				rangeValue = string(metadataJSON)
			}
			queryParts = append(queryParts, k+"::jsonb @> '{\""+rangeKey+"\":"+rangeValue.(string)+"}'::jsonb")
		} else {
			queryParts = append(queryParts, "JSON_EXTRACT("+k+", '$."+rangeKey+"') = '"+escapeDBString(rangeValue.(string))+"'")
		}
	}

	if len(queryParts) == 0 {
		return ""
	}
	query := queryParts[0]
	if len(queryParts) > 1 {
		query = "(" + strings.Join(queryParts, " AND ") + ")"
	}

	return query
}

// whereSlice generates the where slice
func (builder *whereBuilder) whereSlice(k string, v interface{}) string {
	engine := builder.client.Engine()
	if engine == MySQL {
		return "JSON_CONTAINS(" + k + ", CAST('[\"" + v.(string) + "\"]' AS JSON))"
	} else if engine == PostgreSQL {
		return k + "::jsonb @> '[\"" + v.(string) + "\"]'"
	}
	return "EXISTS (SELECT 1 FROM json_each(" + k + ") WHERE value = \"" + v.(string) + "\")"
}
