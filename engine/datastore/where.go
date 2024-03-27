package datastore

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	customtypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"gorm.io/gorm"
)

// CustomWhereInterface is an interface for the CustomWhere clauses
type CustomWhereInterface interface {
	Where(query interface{}, args ...interface{})
	getGormTx() *gorm.DB
}

// CustomWhere add conditions
func (c *Client) CustomWhere(tx CustomWhereInterface, conditions map[string]interface{}, engine Engine) interface{} {
	// Empty accumulator
	varNum := 0

	// Process the conditions
	processConditions(c, tx, conditions, engine, &varNum, nil)

	// Return the GORM tx
	return tx.getGormTx()
}

// txAccumulator is the accumulator struct
type txAccumulator struct {
	CustomWhereInterface
	WhereClauses []string
	Vars         map[string]interface{}
}

// Where is our custom where method
func (tx *txAccumulator) Where(query interface{}, args ...interface{}) {
	tx.WhereClauses = append(tx.WhereClauses, query.(string))

	if len(args) > 0 {
		for _, variables := range args {
			for key, value := range variables.(map[string]interface{}) {
				tx.Vars[key] = value
			}
		}
	}
}

// getGormTx will get the GORM tx
func (tx *txAccumulator) getGormTx() *gorm.DB {
	return nil
}

// processConditions will process all conditions
func processConditions(client ClientInterface, tx CustomWhereInterface, conditions map[string]interface{},
	engine Engine, varNum *int, parentKey *string,
) map[string]interface{} { //nolint:nolintlint,unparam // ignore for now

	for key, condition := range conditions {
		if key == conditionAnd {
			processWhereAnd(client, tx, condition, engine, varNum)
		} else if key == conditionOr {
			processWhereOr(client, tx, conditions[conditionOr], engine, varNum)
		} else if key == conditionGreaterThan {
			varName := "var" + strconv.Itoa(*varNum)
			tx.Where(*parentKey+" > @"+varName, map[string]interface{}{varName: formatCondition(condition, engine)})
			*varNum++
		} else if key == conditionLessThan {
			varName := "var" + strconv.Itoa(*varNum)
			tx.Where(*parentKey+" < @"+varName, map[string]interface{}{varName: formatCondition(condition, engine)})
			*varNum++
		} else if key == conditionGreaterThanOrEqual {
			varName := "var" + strconv.Itoa(*varNum)
			tx.Where(*parentKey+" >= @"+varName, map[string]interface{}{varName: formatCondition(condition, engine)})
			*varNum++
		} else if key == conditionLessThanOrEqual {
			varName := "var" + strconv.Itoa(*varNum)
			tx.Where(*parentKey+" <= @"+varName, map[string]interface{}{varName: formatCondition(condition, engine)})
			*varNum++
		} else if key == conditionExists {
			if condition.(bool) {
				tx.Where(*parentKey + " IS NOT NULL")
			} else {
				tx.Where(*parentKey + " IS NULL")
			}
		} else if StringInSlice(key, client.GetArrayFields()) {
			tx.Where(whereSlice(engine, key, formatCondition(condition, engine)))
		} else if StringInSlice(key, client.GetObjectFields()) {
			tx.Where(whereObject(engine, key, formatCondition(condition, engine)))
		} else {
			if condition == nil {
				tx.Where(key + " IS NULL")
			} else {
				v := reflect.ValueOf(condition)
				switch v.Kind() { //nolint:exhaustive // not all cases are needed
				case reflect.Map:
					if _, ok := condition.(map[string]interface{}); ok {
						processConditions(client, tx, condition.(map[string]interface{}), engine, varNum, &key) //nolint:scopelint // ignore for now
					} else {
						c, _ := json.Marshal(condition) //nolint:errchkjson // this check might break the current code
						var cc map[string]interface{}
						_ = json.Unmarshal(c, &cc)
						processConditions(client, tx, cc, engine, varNum, &key) //nolint:scopelint // ignore for now
					}
				default:
					varName := "var" + strconv.Itoa(*varNum)
					tx.Where(key+" = @"+varName, map[string]interface{}{varName: formatCondition(condition, engine)})
					*varNum++
				}
			}
		}
	}

	return conditions
}

// formatCondition will format the conditions
func formatCondition(condition interface{}, engine Engine) interface{} {
	switch v := condition.(type) {
	case customtypes.NullTime:
		if v.Valid {
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
func processWhereAnd(client ClientInterface, tx CustomWhereInterface, condition interface{}, engine Engine, varNum *int) {
	accumulator := &txAccumulator{
		WhereClauses: make([]string, 0),
		Vars:         make(map[string]interface{}),
	}
	for _, c := range condition.([]map[string]interface{}) {
		processConditions(client, accumulator, c, engine, varNum, nil)
	}

	if len(accumulator.Vars) > 0 {
		tx.Where(" ( "+strings.Join(accumulator.WhereClauses, " AND ")+" ) ", accumulator.Vars)
	} else {
		tx.Where(" ( " + strings.Join(accumulator.WhereClauses, " AND ") + " ) ")
	}
}

// processWhereOr will process the OR statements
func processWhereOr(client ClientInterface, tx CustomWhereInterface, condition interface{}, engine Engine, varNum *int) {
	or := make([]string, 0)
	orVars := make(map[string]interface{})
	for _, cond := range condition.([]map[string]interface{}) {
		statement := make([]string, 0)
		accumulator := &txAccumulator{
			WhereClauses: make([]string, 0),
			Vars:         make(map[string]interface{}),
		}
		processConditions(client, accumulator, cond, engine, varNum, nil)
		statement = append(statement, accumulator.WhereClauses...)
		for varName, varValue := range accumulator.Vars {
			orVars[varName] = varValue
		}
		or = append(or, strings.Join(statement[:], " AND "))
	}

	if len(orVars) > 0 {
		tx.Where(" ( ("+strings.Join(or, ") OR (")+") ) ", orVars)
	} else {
		tx.Where(" ( (" + strings.Join(or, ") OR (") + ") ) ")
	}
}

// escapeDBString will escape the database string
func escapeDBString(s string) string {
	rs := strings.Replace(s, "'", "\\'", -1)
	return strings.Replace(rs, "\"", "\\\"", -1)
}

// whereObject generates the where object
func whereObject(engine Engine, k string, v interface{}) string {
	queryParts := make([]string, 0)

	// we don't know the type, we handle the rangeValue as a map[string]interface{}
	vJSON, _ := json.Marshal(v) //nolint:errchkjson // this check might break the current code

	var rangeV map[string]interface{}
	_ = json.Unmarshal(vJSON, &rangeV)

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
func whereSlice(engine Engine, k string, v interface{}) string {
	if engine == MySQL {
		return "JSON_CONTAINS(" + k + ", CAST('[\"" + v.(string) + "\"]' AS JSON))"
	} else if engine == PostgreSQL {
		return k + "::jsonb @> '[\"" + v.(string) + "\"]'"
	}
	return "EXISTS (SELECT 1 FROM json_each(" + k + ") WHERE value = \"" + v.(string) + "\")"
}
