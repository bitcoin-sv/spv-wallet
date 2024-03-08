package engine

import (
	"database/sql/driver"

	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

type ContactStatus string

const (
	strContactStatusNotConfirmed = "unconfirmed"
	strContactStatusAwaiting     = "awaiting"
	strContactStatusConfirmed    = "confirmed"

	ContactNotConfirmed ContactStatus = strContactStatusNotConfirmed
	ContactAwaitAccept  ContactStatus = strContactStatusAwaiting
	ContactConfirmed    ContactStatus = strContactStatusConfirmed
)

// Scan will scan the value into Struct, implements sql.Scanner interface
func (t *ContactStatus) Scan(value interface{}) error {
	stringValue, err := utils.StrOrBytesToString(value)
	if err != nil {
		return nil
	}

	switch stringValue {
	case strContactStatusNotConfirmed:
		*t = ContactNotConfirmed
	case strContactStatusAwaiting:
		*t = ContactAwaitAccept
	case strContactStatusConfirmed:
		*t = ContactConfirmed
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
