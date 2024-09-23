package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Is checks if the target error is the same as the current error
func (e SPVError) Is(target error) bool {
	t, ok := target.(SPVError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// ErrorNode holds information of parent error and its causes
type ErrorNode struct {
	Err    error        `json:"-"`
	Msg    string       `json:"msg"`
	Type   string       `json:"type"`
	Causes []*ErrorNode `json:"causes,omitempty"`
}

// ToJSON returns the JSON representation of the error node
func (en *ErrorNode) ToJSON() []byte {
	marshalled, err := json.Marshal(en)
	if err != nil {
		return []byte(fmt.Sprintf(`{"err": "%s"}`, en.Msg))
	}
	return marshalled
}

// ToString returns the string representation of the error node
func (en *ErrorNode) ToString() string {
	msg := fmt.Sprintf("'%s' <of type [%s]>", en.Msg, en.Type)
	if len(en.Causes) != 0 {
		strCauses := make([]string, len(en.Causes))
		for i, cause := range en.Causes {
			strCauses[i] = cause.ToString()
		}
		msg += fmt.Sprintf(" was caused by { %s }", strings.Join(strCauses, " AND "))
	}
	return msg
}

// InitialCause returns the initial cause of the error
// Doesn't support joined errors because it's not possible to determine the ONE initial cause (in this case would be multiple)
func (en *ErrorNode) InitialCause() *ErrorNode {
	if len(en.Causes) == 0 {
		return en
	}
	return en.Causes[0].InitialCause()
}

// UnfoldError recursively unfolds the error and its causes
func UnfoldError(err error) *ErrorNode {
	if err == nil {
		return nil
	}

	node := &ErrorNode{
		Err:    err,
		Msg:    err.Error(),
		Type:   reflect.TypeOf(err).String(),
		Causes: []*ErrorNode{},
	}

	// Check for joined errors first (Go 1.20+)
	if jErrors, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range jErrors.Unwrap() {
			node.Causes = append(node.Causes, UnfoldError(e))
		}
	} else if unwrapped := errors.Unwrap(err); unwrapped != nil {
		node.Causes = append(node.Causes, UnfoldError(unwrapped))
	}

	return node
}
