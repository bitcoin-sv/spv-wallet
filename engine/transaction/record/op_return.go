package record

import (
	"github.com/bitcoin-sv/go-sdk/script"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
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
