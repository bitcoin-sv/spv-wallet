package outputs

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/bitcoin-sv/go-sdk/script"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
)

// OpReturn represents an OP_RETURN output specification.
type OpReturn opreturn.Output

func (o *OpReturn) evaluate(context.Context) (annotatedOutputs, error) {
	if len(o.Data) == 0 {
		return nil, txerrors.ErrDraftOpReturnDataRequired
	}

	data, err := o.getData()
	if err != nil {
		return nil, err
	}
	output, err := sdk.CreateOpReturnOutput(data)
	if err != nil {
		if errors.Is(err, script.ErrPartTooBig) {
			return nil, txerrors.ErrDraftOpReturnDataTooLarge
		}
		return nil, spverrors.Wrapf(err, "failed to create OP_RETURN output")
	}

	annotation := transaction.NewDataOutputAnnotation()
	return singleAnnotatedOutput(output, annotation), nil
}

func (o *OpReturn) getData() ([][]byte, error) {
	data := make([][]byte, len(o.Data))
	for i, dataToStore := range o.Data {
		bytes, err := toBytes(dataToStore, o.DataType)
		if err != nil {
			return nil, err
		}
		data[i] = bytes
	}
	return data, nil
}

func toBytes(data string, dataType opreturn.DataType) ([]byte, error) {
	switch dataType {
	case opreturn.DataTypeDefault, opreturn.DataTypeStrings:
		return []byte(data), nil
	case opreturn.DataTypeHexes:
		dataHex, err := hex.DecodeString(data)
		if err != nil {
			return nil, spverrors.Wrapf(err, "failed to decode hex")
		}
		return dataHex, nil
	default:
		return nil, spverrors.Newf("unsupported data type")
	}
}
