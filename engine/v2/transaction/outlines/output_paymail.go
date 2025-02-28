package outlines

import (
	"errors"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-sdk/script"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	pmerrors "github.com/bitcoin-sv/spv-wallet/engine/paymail/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/optional"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

// Paymail represents a paymail output
type Paymail struct {
	To       string
	Satoshis bsv.Satoshis
	Splits   uint64
	From     optional.Param[string]
}

func (p *Paymail) evaluate(ctx *evaluationContext) (annotatedOutputs, error) {
	p.ensureDefaults()

	paymailClient := ctx.Paymail()

	receiverAddress, err := paymailClient.GetSanitizedPaymail(p.To)
	if err != nil {
		return nil, txerrors.ErrReceiverPaymailAddressIsInvalid.Wrap(err)
	}

	err = p.validate()
	if err != nil {
		return nil, err
	}

	sender, err := p.sender(ctx)
	if err != nil {
		return nil, err
	}

	destinations, err := paymailClient.GetP2PDestinations(ctx, receiverAddress, p.Satoshis)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get P2P destinations for paymail %s", p.To)
	}

	if p.wantToSplit() {
		result, err := p.mapToSplitOutputs(destinations, sender)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	result, err := p.mapToAnnotatedOutputs(destinations, sender)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *Paymail) ensureDefaults() {
	if p.Splits == 0 {
		p.Splits = 1
	}
}

func (p *Paymail) validate() error {
	err := p.validateSatoshis()
	if err != nil {
		return err
	}
	err = p.validateSplits()
	if err != nil {
		return err
	}

	return nil
}

func (p *Paymail) validateSatoshis() error {
	if p.Satoshis == 0 {
		return txerrors.ErrOutputValueTooLow
	}
	return nil
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

func (p *Paymail) sender(ctx *evaluationContext) (string, error) {
	if p.From == nil {
		return p.defaultSenderAddress(ctx)
	}

	err := p.validateProvidedSenderPaymail(ctx)
	if err != nil {
		return "", err
	}

	return *p.From, nil
}

func (p *Paymail) validateSplits() error {
	if p.Splits > 1 && uint64(p.Satoshis)%p.Splits != 0 {
		return txerrors.ErrTxOutlinePaymailSatoshisMustBeDivisibleBySplits
	}

	return nil
}

func (p *Paymail) validateProvidedSenderPaymail(ctx *evaluationContext) error {
	var sender = *p.From
	_, err := ctx.Paymail().GetSanitizedPaymail(sender)
	if err != nil {
		return txerrors.ErrSenderPaymailAddressIsInvalid.Wrap(err)
	}
	ownsPaymail, err := ctx.PaymailAddressService().HasPaymailAddress(ctx, ctx.UserID(), sender)
	if errors.Is(err, spverrors.ErrCouldNotFindPaymail) {
		return txerrors.ErrSenderPaymailAddressIsInvalid.Wrap(err)
	}
	if err != nil {
		return spverrors.Wrapf(err, "failed to check if paymail %s belongs to user %s", sender, ctx.UserID())
	}

	if !ownsPaymail {
		return txerrors.ErrSenderPaymailAddressIsInvalid
	}

	return nil
}

func (p *Paymail) defaultSenderAddress(ctx *evaluationContext) (string, error) {
	sender, err := ctx.PaymailAddressService().GetDefaultPaymailAddress(ctx, ctx.UserID())
	if err != nil {
		return "", txerrors.ErrTxOutlineSenderPaymailAddressNoDefault.Wrap(err)
	}
	return sender, nil
}

func (p *Paymail) wantToSplit() bool {
	return p.Splits > 1
}

func (p *Paymail) mapToAnnotatedOutputs(destinations *paymail.PaymentDestinationPayload, sender string) (annotatedOutputs, error) {
	var err error
	result := make(annotatedOutputs, len(destinations.Outputs))
	for i, output := range destinations.Outputs {
		result[i], err = p.createBsvPaymailOutput(output, destinations.Reference, sender)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (p *Paymail) mapToSplitOutputs(destinations *paymail.PaymentDestinationPayload, sender string) (annotatedOutputs, error) {
	if len(destinations.Outputs) != 1 {
		return nil, txerrors.ErrTxOutlinePaymailCannotSplitWhenRecipientSplitting
	}

	output := destinations.Outputs[0]
	satoshisPart := uint64(p.Satoshis) / p.Splits

	var err error
	result := make(annotatedOutputs, p.Splits)
	for i := range p.Splits {
		output.Satoshis = satoshisPart
		result[i], err = p.createBsvPaymailOutput(output, destinations.Reference, sender)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
