package customtypes

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// NullTime wrapper around sql.NullTime
type NullTime struct {
	sql.NullTime
}

// IsZero method is called by bson.IsZero in Mongo for type = NullTime
func (x NullTime) IsZero() bool {
	return !x.Valid
}

// MarshalNullTime is used by GraphQL to marshal the value
func MarshalNullTime(x NullTime) graphql.Marshaler {
	if !x.Valid {
		return graphql.Null
	}

	return graphql.MarshalTime(x.Time)
}

// UnmarshalNullTime is used by GraphQL to unmarshal the value
func UnmarshalNullTime(t interface{}) (NullTime, error) {
	if t == nil {
		return NullTime{sql.NullTime{Valid: false}}, nil
	}

	uTime, err := graphql.UnmarshalTime(t)
	if err != nil {
		return NullTime{}, err
	}

	return NullTime{
		// @mrz: had to remove uTime.UnixMicro() > 0 in Valid (issue was golangci-lint typecheck)
		sql.NullTime{
			Time:  uTime,
			Valid: true,
		},
	}, err
}

// MarshalBSONValue method is called by bson.Marshal in Mongo for type = NullTime
func (x *NullTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !x.Valid {
		return bsontype.Null, nil, nil
	}

	valueType, b, err := bson.MarshalValue(x.Time)
	return valueType, b, err
}

// UnmarshalBSONValue method is called by bson.Unmarshal in Mongo for type = NullTime
func (x *NullTime) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	raw := bson.RawValue{Type: t, Value: data}
	uTime := time.Time{}

	if err := raw.Unmarshal(&uTime); err != nil {
		return err
	}

	if raw.Value == nil {
		x.Valid = false
		return nil
	}

	x.Valid = true
	x.Time = uTime
	return nil
}

// MarshalJSON method is called by the JSON marshaller
func (x *NullTime) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}

	b, err := json.Marshal(x.Time)
	return b, err
}

// UnmarshalJSON method is called by the JSON unmarshaller
func (x *NullTime) UnmarshalJSON(data []byte) error {
	x.Valid = false

	if data == nil {
		return nil
	}

	var timeString string
	if err := json.Unmarshal(data, &timeString); err != nil {
		return err
	}
	if timeString == "" {
		return nil
	}

	uTime, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return err
	}

	x.Valid = true
	x.Time = uTime
	return nil
}
