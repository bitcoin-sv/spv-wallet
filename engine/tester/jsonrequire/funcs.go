package jsonrequire

import (
	"fmt"
	"strings"
	"text/template"
)

var funcsMap = template.FuncMap{
	"matchTimestamp":     matchTimestamp,
	"matchURL":           matchURL,
	"orEmpty":            orEmpty,
	"matchID64":          matchID64,
	"matchHexWithLength": matchHexWithLength,
	"matchHex":           matchHex,
	"matchBEEF":          matchBEEF,
	"matchAddress":       matchAddress,
	"matchNumber":        matchNumber,
	"anything":           anything,
	"matchTxByFormat":    matchTxByFormat,
	"matchDestination":   matchDestination,
	"containsAll":        containsAll,
}

func anything() string {
	return `"*"`
}

func matchTimestamp() string {
	return regexPlaceholder(`^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d+(Z|[+-]\\d{2}:\\d{2})$`)
}

func matchURL() string {
	return regexPlaceholder(`^(https?|ftp):\\/\\/[^\\s/$.?#].[^\\s]*$`)
}

func matchID64() string {
	return regexPlaceholder(`^[a-zA-Z0-9]{64}$`)
}

func matchHexWithLength(length int) string {
	return regexPlaceholder(
		fmt.Sprintf(`^[a-fA-F0-9]{%d}$`, length),
	)
}

func matchHex() string {
	return regexPlaceholder(`^[a-fA-F0-9]+$`)
}

func matchBEEF() string {
	return regexPlaceholder(`^0100(beef|BEEF)[a-fA-F0-9]+$`)
}

func matchTxByFormat(format string) string {
	switch strings.ToLower(format) {
	case "beef":
		return matchBEEF()
	case "raw":
		return matchHex()
	default:
		panic(fmt.Sprintf("unsupported tx format: %s", format))
	}
}

func matchDestination() string {
	return regexPlaceholder("^1-destination-.{32}$")
}

// matchAddress returns a regex that matches a bitcoin address
// NOTE: Only P2PKH (mainnet) addresses are supported
func matchAddress() string {
	return regexPlaceholder(`^1[a-km-zA-HJ-NP-Z1-9]{24,33}$`)
}

func matchNumber() string {
	return regexPlaceholder(`^\\d+$`)
}

func containsAll(parts []string) string {
	if len(parts) == 0 {
		return "*"
	}
	partsRegex := strings.Builder{}
	partsRegex.WriteString("^")
	for _, part := range parts {
		partsRegex.WriteString(fmt.Sprintf(`(.*%s)`, part))
	}
	partsRegex.WriteString(".*$")
	s := partsRegex.String()
	return regexPlaceholder(s)
}

// regexPlaceholder adds slashes at the beginning and end of a string
// it is an indicator for placeholder matcher algorithm that it is a regex
func regexPlaceholder(statement string) string {
	return fmt.Sprintf("/%s/", statement)
}

func orEmpty(statement string) string {
	if !containsRegex(statement) {
		return statement
	}

	regex := extractRegex(statement)

	return regexPlaceholder(
		fmt.Sprintf(`(%s)|^$`, regex),
	)
}
