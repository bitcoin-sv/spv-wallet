package gqlgen

/*
// MarshalNullString is ...
func MarshalNullString(ns null.String) graphql.Marshaler {
	if !ns.Valid {
		// this is also important, so we can detect if this scalar is used in a not null context and return an appropriate error
		return graphql.Null
	}
	return graphql.MarshalString(ns.String)
}

// UnmarshalNullString is ...
func UnmarshalNullString(v interface{}) (null.String, error) {
	if v == nil {
		return null.String{Valid: false}, nil
	}
	// again you can delegate to the default implementation to save yourself some work.
	s, err := graphql.UnmarshalString(v)
	return null.String{String: s, Valid: true}, err
}

// MarshalNullBool is ...
func MarshalNullBool(ns null.Bool) graphql.Marshaler {
	if !ns.Valid {
		// this is also important, so we can detect if this scalar is used in a not null context and return an appropriate error
		return graphql.Null
	}
	return graphql.MarshalBoolean(ns.Bool)
}

// UnmarshalNullBool is ...
func UnmarshalNullBool(v interface{}) (null.Bool, error) {
	if v == nil {
		return null.Bool{Valid: false}, nil
	}
	// again you can delegate to the default implementation to save yourself some work.
	s, err := graphql.UnmarshalBoolean(v)
	return null.Bool{Bool: s, Valid: true}, err
}

// MarshalNullTime is ...
func MarshalNullTime(ns null.Time) graphql.Marshaler {
	if !ns.Valid {
		// this is also important, so we can detect if this scalar is used in a not null context and return an appropriate error
		return graphql.Null
	}
	return graphql.MarshalTime(ns.Time)
}

// UnmarshalNullTime is ...
func UnmarshalNullTime(v interface{}) (null.Time, error) {
	if v == nil {
		return null.Time{Valid: false}, nil
	}
	// again you can delegate to the default implementation to save yourself some work.
	s, err := graphql.UnmarshalTime(v)
	return null.Time{Time: s, Valid: true}, err
}
*/
