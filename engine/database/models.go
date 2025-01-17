package database

// Models returns a list of all models, e.g. for migrations.
func Models() []any {
	return []any{
		TrackedTransaction{},
		TrackedOutput{},
		Data{},
		User{},
		Paymail{},
		Address{},
		UsersUTXO{},
		Operation{},
	}
}
