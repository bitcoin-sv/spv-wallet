package outputs

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-sdk/script"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	pmerrors "github.com/bitcoin-sv/spv-wallet/engine/paymail/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/evaluation"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	paymailreq "github.com/bitcoin-sv/spv-wallet/models/request/paymail"
)

// Paymail represents a paymail output
type Paymail paymailreq.Output

func (p *Paymail) evaluate(ctx evaluation.Context) (annotatedOutputs, error) {
	paymailClient := ctx.Paymail()

	paymailAddress, err := paymailClient.GetSanitizedPaymail(p.To)
	if err != nil {
		return nil, spverrors.ErrPaymailAddressIsInvalid.Wrap(err)
	}

	if p.Satoshis == 0 {
		return nil, txerrors.ErrOutputValueTooLow
	}

	destinations, err := paymailClient.GetP2PDestinations(ctx, paymailAddress, p.Satoshis)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get P2P destinations for paymail %s", p.To)
	}

	result := make(annotatedOutputs, len(destinations.Outputs))
	for i, output := range destinations.Outputs {
		result[i], err = p.createBsvPaymailOutput(output, destinations.Reference)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (p *Paymail) createBsvPaymailOutput(output *paymail.PaymentOutput, reference string) (*annotatedOutput, error) {
	lockingScript, err := script.NewFromHex(output.Script)
	if err != nil {
		return nil, pmerrors.ErrPaymailHostInvalidResponse.Wrap(err)
	}

	return &annotatedOutput{
		TransactionOutput: &sdk.TransactionOutput{
			Satoshis:      output.Satoshis,
			LockingScript: lockingScript,
		},
		OutputAnnotation: &transaction.OutputAnnotation{
			Bucket: transaction.BucketBSV,
			Paymail: &transaction.PaymailAnnotation{
				Receiver:  p.To,
				Reference: reference,
			},
		},
	}, nil
}
