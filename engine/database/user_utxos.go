package database

import "time"

// UserUtxos is a table holding user's Unspent Transaction Outputs (UTXOs).
type UserUtxos struct {
	UserID                       string `gorm:"primaryKey;uniqueIndex:idx_window,sort:asc,priority:1"`
	TxID                         string `gorm:"primaryKey;uniqueIndex:idx_window,sort:asc,priority:4"`
	Vout                         uint32 `gorm:"primaryKey;uniqueIndex:idx_window,sort:asc,priority:5"`
	Satoshis                     uint64
	UnlockingScriptEstimatedSize uint64
	Bucket                       string    `gorm:"check:chk_not_data_bucket,bucket <> 'data'"`
	CreatedAt                    time.Time `gorm:"uniqueIndex:idx_window,sort:asc,priority:3"`
	TouchedAt                    time.Time `gorm:"uniqueIndex:idx_window,sort:asc,priority:2"`
}
