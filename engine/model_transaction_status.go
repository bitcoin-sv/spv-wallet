package engine

import (
	"database/sql/driver"

	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// DraftStatus draft transaction status
type DraftStatus string

const (
	// DraftStatusDraft is when the transaction is a draft
	DraftStatusDraft DraftStatus = statusDraft

	// DraftStatusCanceled is when the draft is canceled
	DraftStatusCanceled DraftStatus = statusCanceled

	// DraftStatusExpired is when the draft has expired
	DraftStatusExpired DraftStatus = statusExpired

	// DraftStatusComplete is when the draft transaction is complete
	DraftStatusComplete DraftStatus = statusComplete
)

// Scan will scan the value into Struct, implements sql.Scanner interface
func (t *DraftStatus) Scan(value interface{}) error {
	stringValue, err := utils.StrOrBytesToString(value)
	if err != nil {
		return nil
	}

	switch stringValue {
	case statusDraft:
		*t = DraftStatusDraft
	case statusCanceled:
		*t = DraftStatusCanceled
	case statusExpired:
		*t = DraftStatusExpired
	case statusComplete:
		*t = DraftStatusComplete
	}

	return nil
}

// Value return json value, implement driver.Valuer interface
func (t DraftStatus) Value() (driver.Value, error) {
	return string(t), nil
}
