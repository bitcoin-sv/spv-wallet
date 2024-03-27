package datastore

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	customtypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type CustomWhereInterface interface {
	Where(query interface{}, args ...interface{}) *gorm.DB
}

type txAccumulator struct {
	CustomWhereInterface
	WhereClauses []string
	Vars         map[string]interface{}
}

// Where is our custom where method
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

// WhereBuilder holds a state during custom where preparation
type WhereBuilder struct {
	client ClientInterface
	gdb    *gorm.DB
	varNum int
}

// ApplyCustomWhere adds conditions to the gorm db instance
func ApplyCustomWhere(client ClientInterface, gdb *gorm.DB, conditions map[string]interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error processing conditions: %v", r)
		}
	}()

	builder := &WhereBuilder{
		client: client,
		gdb:    gdb,
		varNum: 0,
	}

	builder.processConditions(gdb, conditions, nil)
	return nil
}

func (builder *WhereBuilder) nextVarName() string {
	varName := "var" + strconv.Itoa(builder.varNum)
	builder.varNum++
	return varName
}

func getColumnName(columnName string, model interface{}) string {
	sch, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		panic(fmt.Errorf("cannot parse a model %v", model))
	}
	if field, ok := sch.FieldsByDBName[columnName]; ok {
		return field.DBName
	}

	if field, ok := sch.FieldsByName[columnName]; ok {
		return field.DBName
	}

	panic(fmt.Errorf("column %s does not exist in the model", columnName))
}

func (builder *WhereBuilder) applyCondition(tx CustomWhereInterface, key string, operator string, condition interface{}) {
	columnName := getColumnName(key, builder.gdb.Statement.Model)

	varName := builder.nextVarName()
	query := fmt.Sprintf("%s %s @%s", columnName, operator, varName)
	tx.Where(query, map[string]interface{}{varName: builder.formatCondition(condition)})
}

func (builder *WhereBuilder) applyExistsCondition(tx CustomWhereInterface, key string, condition bool) {
	columnName := getColumnName(key, builder.gdb.Statement.Model)

	operator := "IS NULL"
	if condition {
		operator = "IS NOT NULL"
	}
	tx.Where(columnName + " " + operator)
}

// processConditions will process all conditions
func (builder *WhereBuilder) processConditions(tx CustomWhereInterface, conditions map[string]interface{}, parentKey *string,
) {
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
func (builder *WhereBuilder) formatCondition(condition interface{}) interface{} {
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
func (builder *WhereBuilder) processWhereAnd(tx CustomWhereInterface, condition interface{}) {
	accumulator := &txAccumulator{
		WhereClauses: make([]string, 0),
		Vars:         make(map[string]interface{}),
	}
	for _, c := range condition.([]map[string]interface{}) {
		builder.processConditions(accumulator, c, nil)
	}

	tx.Where(" ( "+strings.Join(accumulator.WhereClauses, " AND ")+" ) ", accumulator.Vars)
}

// processWhereOr will process the OR statements
func (builder *WhereBuilder) processWhereOr(tx CustomWhereInterface, condition interface{}) {
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

	tx.Where(" ( ("+strings.Join(or, ") OR (")+") ) ", orVars)
}

// escapeDBString will escape the database string
func escapeDBString(s string) string {
	rs := strings.Replace(s, "'", "\\'", -1)
	return strings.Replace(rs, "\"", "\\\"", -1)
}

// whereObject generates the where object
func (builder *WhereBuilder) whereObject(k string, v interface{}) string {
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
func (builder *WhereBuilder) whereSlice(k string, v interface{}) string {
	engine := builder.client.Engine()
	if engine == MySQL {
		return "JSON_CONTAINS(" + k + ", CAST('[\"" + v.(string) + "\"]' AS JSON))"
	} else if engine == PostgreSQL {
		return k + "::jsonb @> '[\"" + v.(string) + "\"]'"
	}
	return "EXISTS (SELECT 1 FROM json_each(" + k + ") WHERE value = \"" + v.(string) + "\")"
}
