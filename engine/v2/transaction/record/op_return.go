package record

import (
	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/txmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

func getDataFromOpReturn(lockingScript *script.Script) ([]byte, error) {
	if !lockingScript.IsData() {
		return nil, txerrors.ErrAnnotationMismatch
	}

	chunks, err := lockingScript.Chunks()
	if err != nil {
		return nil, txerrors.ErrParsingScript.Wrap(err)
	}

	startIndex := 2
	if chunks[0].Op == script.OpRETURN {
		startIndex = 1
	}

	var bytes []byte
	for _, chunk := range chunks[startIndex:] {
		if chunk.Op > script.OpPUSHDATA4 || chunk.Op == script.OpZERO {
			return nil, txerrors.ErrOnlyPushDataAllowed
		}
		bytes = append(bytes, chunk.Data...)
	}

	return bytes, nil
}

func processDataOutputs(tx *trx.Transaction, userID string, annotations *transaction.Annotations) ([]txmodels.NewOutput, error) {
	txID := tx.TxID().String()

	var err error
	var dataOutputs []txmodels.NewOutput //nolint: prealloc

	for vout, annotation := range annotations.Outputs {
		if vout >= len(tx.Outputs) {
			return nil, txerrors.ErrAnnotationIndexOutOfRange
		}
		outpoint := bsv.Outpoint{TxID: txID}
		outpoint.Vout, err = conv.IntToUint32(vout)
		if err != nil {
			return nil, txerrors.ErrAnnotationIndexConversion.Wrap(err)
		}

		lockingScript := tx.Outputs[vout].LockingScript

		if annotation.Bucket != bucket.Data {
			continue
		}

		data, err := getDataFromOpReturn(lockingScript)
		if err != nil {
			return nil, err
		}
		dataOutputs = append(dataOutputs, txmodels.NewOutputForData(outpoint, userID, data))
	}

	return dataOutputs, nil
}
