package record

import (
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"iter"

	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
)

type operationWrapper struct {
	entity *database.Operation
}

func (w *operationWrapper) add(satoshi bsv.Satoshis) {
	signedSatoshi, err := conv.Uint64ToInt64(uint64(satoshi))
	if err != nil {
		panic(err)
	}
	w.entity.Value = w.entity.Value + signedSatoshi
}

func (w *operationWrapper) subtract(satoshi bsv.Satoshis) {
	signedSatoshi, err := conv.Uint64ToInt64(uint64(satoshi))
	if err != nil {
		panic(err)
	}
	w.entity.Value = w.entity.Value - signedSatoshi
}

func toOperationEntities(wrappers iter.Seq[*operationWrapper]) iter.Seq[*database.Operation] {
	return func(yield func(*database.Operation) bool) {
		for wrapper := range wrappers {
			yield(wrapper.entity)
		}
	}
}
