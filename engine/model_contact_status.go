package engine

import (
	"database/sql/driver"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// ContactStatus represents statuses of contact model.
type ContactStatus string

const (
	// ContactNotConfirmed is a status telling that the contact model as not confirmed yet.
	ContactNotConfirmed ContactStatus = "unconfirmed"
	// ContactAwaitAccept is a status telling that the contact model as invitation to add to contacts.
	ContactAwaitAccept ContactStatus = "awaiting"
	// ContactConfirmed is a status telling that the contact model as confirmed.
	ContactConfirmed ContactStatus = "confirmed"
	// ContactRejected is a status telling that the contact invitation was rejected by user.
	ContactRejected ContactStatus = "rejected"
)

var contactStatusMapper = NewEnumStringMapper(ContactNotConfirmed, ContactAwaitAccept, ContactConfirmed, ContactRejected)

// Scan will scan the value into Struct, implements sql.Scanner interface
func (t *ContactStatus) Scan(value interface{}) error {
	stringValue, err := utils.StrOrBytesToString(value)
	if err != nil {
		return nil
	}

	status, ok := contactStatusMapper.Get(stringValue)
	if !ok {
		return spverrors.Newf("invalid contact status: %s", stringValue)
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
