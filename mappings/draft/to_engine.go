package mappingsdraft

import (
	"errors"
	"reflect"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/outputs"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
	paymailreq "github.com/bitcoin-sv/spv-wallet/models/request/paymail"
	"github.com/mitchellh/mapstructure"
)

// ToEngine converts a draft transaction request model to the engine model.
func ToEngine(xPubID string, tx *request.DraftTransaction) (*draft.TransactionSpec, error) {
	spec := &draft.TransactionSpec{
		XPubID: xPubID,
	}
	config := mapstructure.DecoderConfig{
		DecodeHook: outputsHookFunc(),
		Result:     &spec,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return nil, spverrors.Wrapf(err, spverrors.ErrCannotMapFromModel.Error())
	}

	err = decoder.Decode(tx)
	if err != nil {
		return nil, spverrors.Wrapf(err, spverrors.ErrCannotMapFromModel.Error())
	}

	return spec, nil
}

func outputsHookFunc() mapstructure.DecodeHookFunc {
	return func(_ reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		specs := outputs.NewSpecifications()
		reqOutputs, ok := data.([]request.Output)
		if !ok {
			return data, nil
		}
		if to != reflect.TypeOf(specs) {
			return data, nil
		}

		for _, out := range reqOutputs {
			spec, err := outputSpecFromRequest(out)
			if err != nil {
				return nil, err
			}
			specs.Add(spec)
		}
		return specs, nil
	}
}

func outputSpecFromRequest(req request.Output) (outputs.Spec, error) {
	switch o := req.(type) {
	case opreturn.Output:
		out := outputs.OpReturn(o)
		return &out, nil
	case paymailreq.Output:
		out := outputs.Paymail(o)
		return &out, nil
	default:
		return nil, errors.New("unsupported output type")
	}
}
