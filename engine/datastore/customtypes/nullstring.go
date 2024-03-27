package customtypes

import (
	"database/sql"
	"encoding/json"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// NullString wrapper around sql.NullString
type NullString struct {
	sql.NullString
}

// IsZero method is called by bson.IsZero in Mongo for type = NullTime
func (x NullString) IsZero() bool {
	return !x.Valid
}

// MarshalNullString is used by graphql to marshal into a string
func MarshalNullString(x NullString) graphql.Marshaler {
	if !x.Valid {
		return graphql.Null
	}

	return graphql.MarshalString(x.String)
}

// UnmarshalNullString is used by graphql to unmarshal from a NullString into a string
func UnmarshalNullString(s interface{}) (NullString, error) {
	if s == nil {
		return NullString{sql.NullString{Valid: false}}, nil
	}

	uString, err := graphql.UnmarshalString(s)
	if err != nil {
		return NullString{}, err
	}

	return NullString{
		sql.NullString{
			String: uString,
			Valid:  true,
		},
	}, err
}

// MarshalBSONValue method is called by bson.Marshal in Mongo for type = NullString
func (x *NullString) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !x.Valid {
		return bsontype.Null, nil, nil
	}

	valueType, b, err := bson.MarshalValue(x.String)
	return valueType, b, err
}

// UnmarshalBSONValue method is called by bson.Unmarshal in Mongo for type = NullString
func (x *NullString) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	raw := bson.RawValue{Type: t, Value: data}

	var uString string
	if err := raw.Unmarshal(&uString); err != nil {
		return err
	}

	if raw.Value == nil {
		x.Valid = false
		return nil
	}

	x.Valid = true
	x.String = uString
	return nil
}

// MarshalJSON method is called by the JSON marshaller
func (x *NullString) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}

	b, err := json.Marshal(x.String)
	return b, err
}

// UnmarshalJSON method is called by the JSON unmarshaller
func (x *NullString) UnmarshalJSON(data []byte) error {
	x.Valid = false

	if data == nil {
		return nil
	}

	var nullString string
	if err := json.Unmarshal(data, &nullString); err != nil {
		return err
	}

	x.Valid = true
	x.String = nullString
	return nil
}
