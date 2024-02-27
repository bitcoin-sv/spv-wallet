package engine

import (
	"database/sql/driver"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

type ContactStatus string

const (
	ContactStatusNotConf     = notConfirmed
	ContactStatusAwaitAccept = awaitingAcceptance
	ContactStatusConfirmed   = confirmed
)

// Scan will scan the value into Struct, implements sql.Scanner interface
func (t *ContactStatus) Scan(value interface{}) error {
	stringValue, err := utils.StrOrBytesToString(value)
	if err != nil {
		return nil
	}

	switch stringValue {
	case notConfirmed:
		*t = ContactStatusNotConf
	case awaitingAcceptance:
		*t = ContactStatusAwaitAccept
	case confirmed:
		*t = ContactStatusConfirmed
	}

	return nil
}

// Value return json value, implement driver.Valuer interface
func (t ContactStatus) Value() (driver.Value, error) {
	return string(t), nil
}

// String is the string version of the status
func (t ContactStatus) String() string {
	return string(t)
}
