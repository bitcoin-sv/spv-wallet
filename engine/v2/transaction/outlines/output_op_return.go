package outlines

import (
	"encoding/hex"
	"errors"

	"github.com/bitcoin-sv/go-sdk/script"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
)

// DataType represents the type of data in the OP_RETURN output.
type DataType int

// Enum values for DataType
const (
	// DataTypeStrings marks the data type as strings.
	DataTypeStrings DataType = iota
	// DataTypeHexes marks the data type as hexes.
	DataTypeHexes
)

// OpReturn represents an OP_RETURN output specification.
type OpReturn struct {
	DataType DataType `json:"dataType,omitempty"`
	Data     []string `json:"data"`
}

func (o *OpReturn) evaluate(*evaluationContext) (annotatedOutputs, error) {
	if len(o.Data) == 0 {
		return nil, txerrors.ErrTxOutlineOpReturnDataRequired
	}

	data, err := o.getData()
	if err != nil {
		return nil, err
	}

	output, err := sdk.CreateOpReturnOutput(data)
	if err != nil {
		if errors.Is(err, script.ErrPartTooBig) {
			return nil, txerrors.ErrTxOutlineOpReturnDataTooLarge
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

func toBytes(data string, dataType DataType) ([]byte, error) {
	switch dataType {
	case DataTypeStrings:
		return []byte(data), nil
	case DataTypeHexes:
		dataHex, err := hex.DecodeString(data)
		if err != nil {
			return nil, txerrors.ErrFailedToDecodeHex.Wrap(err)
		}
		return dataHex, nil
	default:
		return nil, txerrors.ErrTxOutlineOpReturnUnsupportedDataType
	}
}
