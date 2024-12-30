package jsonrequire

import (
	"fmt"
	"text/template"
)

var funcsMap = template.FuncMap{
	"matchTimestamp":     matchTimestamp,
	"matchURL":           matchURL,
	"orEmpty":            orEmpty,
	"matchID64":          matchID64,
	"matchHexWithLength": matchHexWithLength,
}

func matchTimestamp() string {
	return `/^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d+(Z|[+-]\\d{2}:\\d{2})$/`
}

func matchURL() string {
	return `/^(https?|ftp):\\/\\/[^\\s/$.?#].[^\\s]*$/`
}

func matchID64() string {
	return `/^[a-zA-Z0-9]{64}$/`
}

func matchHexWithLength(length int) string {
	return fmt.Sprintf(`/^[a-f0-9]{%d}$/`, length)
}

func orEmpty(statement string) string {
	if !containsRegex(statement) {
		return statement
	}

	regex := extractRegex(statement)

	return fmt.Sprintf(`/(%s)|^$/`, regex)
}
