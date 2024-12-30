package database

// Models returns a list of all models, e.g. for migrations.
func Models() []any {
	return []any{
		TrackedTransaction{},
		Output{},
		Data{},
		User{},
		Paymail{},
	}
}
