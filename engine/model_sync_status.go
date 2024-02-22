package engine

import (
	"database/sql/driver"

	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// SyncStatus sync status
type SyncStatus string

const (
	// SyncStatusPending is when the sync is pending (blocked by other constraints)
	SyncStatusPending SyncStatus = statusPending

	// SyncStatusReady is when the sync is ready (waiting for workers)
	SyncStatusReady SyncStatus = statusReady

	// SyncStatusProcessing is when the sync is processing (worker is running task)
	SyncStatusProcessing SyncStatus = statusProcessing

	// SyncStatusCanceled is when the sync is canceled
	SyncStatusCanceled SyncStatus = statusCanceled

	// SyncStatusSkipped is when the sync is skipped
	SyncStatusSkipped SyncStatus = statusSkipped

	// SyncStatusError is when the sync has an error
	SyncStatusError SyncStatus = statusError

	// SyncStatusComplete is when the sync is complete
	SyncStatusComplete SyncStatus = statusComplete
)

// Scan will scan the value into Struct, implements sql.Scanner interface
func (t *SyncStatus) Scan(value interface{}) error {
	stringValue, err := utils.StrOrBytesToString(value)
	if err != nil {
		return nil
	}

	switch stringValue {
	case statusPending:
		*t = SyncStatusPending
	case statusReady:
		*t = SyncStatusReady
	case statusProcessing:
		*t = SyncStatusProcessing
	case statusCanceled:
		*t = SyncStatusCanceled
	case statusError:
		*t = SyncStatusError
	case statusComplete:
		*t = SyncStatusComplete
	case statusSkipped:
		*t = SyncStatusSkipped
	}

	return nil
}

// Value return json value, implement driver.Valuer interface
func (t SyncStatus) Value() (driver.Value, error) {
	return string(t), nil
}

// String is the string version of the status
func (t SyncStatus) String() string {
	return string(t)
}
