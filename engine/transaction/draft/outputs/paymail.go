package outputs

import (
	"errors"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-sdk/script"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	pmerrors "github.com/bitcoin-sv/spv-wallet/engine/paymail/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/evaluation"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	paymailreq "github.com/bitcoin-sv/spv-wallet/models/request/paymail"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

// Paymail represents a paymail output
type Paymail paymailreq.Output

func (p *Paymail) evaluate(ctx evaluation.Context) (annotatedOutputs, error) {
	paymailClient := ctx.Paymail()

	receiverAddress, err := paymailClient.GetSanitizedPaymail(p.To)
	if err != nil {
		return nil, txerrors.ErrReceiverPaymailAddressIsInvalid.Wrap(err)
	}

	if p.Satoshis == 0 {
		return nil, txerrors.ErrOutputValueTooLow
	}

	sender, err := p.sender(ctx)
	if err != nil {
		return nil, err
	}

	destinations, err := paymailClient.GetP2PDestinations(ctx, receiverAddress, p.Satoshis)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get P2P destinations for paymail %s", p.To)
	}

	result := make(annotatedOutputs, len(destinations.Outputs))
	for i, output := range destinations.Outputs {
		result[i], err = p.createBsvPaymailOutput(output, destinations.Reference, sender)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (p *Paymail) createBsvPaymailOutput(output *paymail.PaymentOutput, reference string, from string) (*annotatedOutput, error) {
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
			Bucket: bucket.BSV,
			Paymail: &transaction.PaymailAnnotation{
				Receiver:  p.To,
				Reference: reference,
				Sender:    from,
			},
		},
	}, nil
}

func (p *Paymail) sender(ctx evaluation.Context) (string, error) {
	if p.From == nil {
		return p.defaultSenderAddress(ctx)
	}

	err := p.validateProvidedSenderPaymail(ctx)
	if err != nil {
		return "", err
	}

	return *p.From, nil
}

func (p *Paymail) validateProvidedSenderPaymail(ctx evaluation.Context) error {
	var sender = *p.From
	_, err := ctx.Paymail().GetSanitizedPaymail(sender)
	if err != nil {
		return txerrors.ErrSenderPaymailAddressIsInvalid.Wrap(err)
	}
	ownsPaymail, err := ctx.PaymailAddressService().HasPaymailAddress(ctx, ctx.XPubID(), sender)
	if errors.Is(err, spverrors.ErrCouldNotFindPaymail) {
		return txerrors.ErrSenderPaymailAddressIsInvalid.Wrap(err)
	}
	if err != nil {
		return spverrors.Wrapf(err, "failed to check if paymail %s belongs to xpub %s", sender, ctx.XPubID())
	}

	if !ownsPaymail {
		return txerrors.ErrSenderPaymailAddressIsInvalid
	}

	return nil
}

func (p *Paymail) defaultSenderAddress(ctx evaluation.Context) (string, error) {
	sender, err := ctx.PaymailAddressService().GetDefaultPaymailAddress(ctx, ctx.XPubID())
	if err != nil {
		return "", txerrors.ErrDraftSenderPaymailAddressNoDefault.Wrap(err)
	}
	return sender, nil
}
