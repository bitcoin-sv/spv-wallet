package engine

import (
	"context"
	"fmt"

	"github.com/mrz1836/go-cachestore"
)

const (
	lockKeyProcessBroadcastTx = "process-broadcast-transaction-%s" // + Tx ID
	lockKeyProcessP2PTx       = "process-p2p-transaction-%s"       // + Tx ID
	lockKeyProcessSyncTx      = "process-sync-transaction-task"
	lockKeyRecordTx           = "action-record-transaction-%s" // + Tx ID
	lockKeyReserveUtxo        = "utxo-reserve-xpub-id-%s"      // + Xpub ID
)

// newWriteLock will take care of creating a lock and defer
func newWriteLock(ctx context.Context, lockKey string, cacheStore cachestore.LockService) (func(), error) {
	secret, err := cacheStore.WriteLock(ctx, lockKey, defaultCacheLockTTL)
	return func() {
		// context is not set, since the req could be canceled, but unlocking should never be stopped
		_, _ = cacheStore.ReleaseLock(context.Background(), lockKey, secret)
	}, err
}

// newWaitWriteLock will take care of creating a lock and defer
func newWaitWriteLock(ctx context.Context, lockKey string, cacheStore cachestore.LockService) (func(), error) {
	secret, err := cacheStore.WaitWriteLock(ctx, lockKey, defaultCacheLockTTL, defaultCacheLockTTW)
	return func() {
		// context is not set, since the req could be canceled, but unlocking should never be stopped
		_, _ = cacheStore.ReleaseLock(context.Background(), lockKey, secret)
	}, err
}

func getWaitWriteLockForPaymail(ctx context.Context, cs cachestore.LockService, id string) (unlock func(), err error) {
	lockKey := fmt.Sprintf("process-paymail-%s", id)
	unlock, err = newWaitWriteLock(ctx, lockKey, cs)
	return
}

func getWaitWriteLockForXpub(ctx context.Context, cs cachestore.LockService, id string) (unlock func(), err error) {
	lockKey := fmt.Sprintf("action-xpub-id-%s", id)
	unlock, err = newWaitWriteLock(ctx, lockKey, cs)
	return
}
