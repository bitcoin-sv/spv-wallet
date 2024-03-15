package engine

import (
	"database/sql/driver"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

type ContactStatus string

const (
	ContactNotConfirmed ContactStatus = "unconfirmed"
	ContactAwaitAccept  ContactStatus = "awaiting"
	ContactConfirmed    ContactStatus = "confirmed"
)

var contactStatusMapper = NewEnumStringMapper(ContactNotConfirmed, ContactAwaitAccept, ContactConfirmed)

// Scan will scan the value into Struct, implements sql.Scanner interface
func (t *ContactStatus) Scan(value interface{}) error {
	stringValue, err := utils.StrOrBytesToString(value)
	if err != nil {
		return nil
	}

	status, ok := contactStatusMapper.Get(stringValue)
	if !ok {
		return fmt.Errorf("invalid contact status: %s", stringValue)
	}
	*t = status

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
