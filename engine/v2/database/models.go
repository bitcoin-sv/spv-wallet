package database

// Models returns a list of all models, e.g. for migrations.
func Models() []any {
	return []any{
		TrackedTransaction{},
		TrackedOutput{},
		TxInput{},
		Data{},
		User{},
		Paymail{},
		Address{},
		UserUTXO{},
		Operation{},
		UserContact{},
	}
}
