package outline

import (
	"errors"
	"reflect"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
	paymailreq "github.com/bitcoin-sv/spv-wallet/models/request/paymail"
	"github.com/mitchellh/mapstructure"
)

// Request is a transaction outline request model.
type Request request.TransactionSpecification

// ToEngine converts a transaction outline request model to the engine model.
func (tx Request) ToEngine(userID string) (*outlines.TransactionSpec, error) {
	spec := &outlines.TransactionSpec{
		UserID: userID,
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
		specs := outlines.NewOutputsSpec()
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

func outputSpecFromRequest(req request.Output) (outlines.OutputSpec, error) {
	switch o := req.(type) {
	case opreturn.Output:
		out := outlines.OpReturn(o)
		return &out, nil
	case paymailreq.Output:
		out := outlines.Paymail(o)
		return &out, nil
	default:
		return nil, errors.New("unsupported output type")
	}
}
