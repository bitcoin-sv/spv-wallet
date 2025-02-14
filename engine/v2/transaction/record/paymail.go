package record

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

type paymailInfo struct {
	annotation       *transaction.PaymailAnnotation
	vouts            map[uint32]struct{}
	skipNotification bool
}

func newPaymailInfo() paymailInfo {
	return paymailInfo{
		annotation:       nil,
		vouts:            map[uint32]struct{}{},
		skipNotification: false,
	}
}

func (pi *paymailInfo) empty() bool {
	return pi.annotation == nil
}

func (pi *paymailInfo) equalsToAnnotation(annotation *transaction.PaymailAnnotation) bool {
	return *pi.annotation == *annotation
}

func (pi *paymailInfo) add(vout int, annotation *transaction.PaymailAnnotation) error {
	if !pi.empty() && !pi.equalsToAnnotation(annotation) {
		return txerrors.ErrMultiPaymailRecipientsNotSupported
	}

	vout32, err := conv.IntToUint32(vout)
	if err != nil {
		return txerrors.ErrAnnotationIndexConversion.Wrap(err)
	}

	pi.annotation = annotation
	pi.vouts[vout32] = struct{}{}
	return nil
}

func (pi *paymailInfo) hasVOut(vout uint32) bool {
	if pi.empty() {
		return false
	}
	_, ok := pi.vouts[vout]
	return ok
}

func (pi *paymailInfo) Sender() string {
	if pi.empty() {
		return ""
	}
	return pi.annotation.Sender
}

func (pi *paymailInfo) Receiver() string {
	if pi.empty() {
		return ""
	}
	return pi.annotation.Receiver
}

func (pi *paymailInfo) Reference() string {
	if pi.empty() {
		return ""
	}
	return pi.annotation.Reference
}

func (f *txFlow) processPaymailOutputs(annotations transaction.Annotations) (paymailInfo, error) {
	info := newPaymailInfo()

	for vout, annotation := range annotations.Outputs {
		if annotation.Paymail == nil {
			continue
		}
		if annotation.Bucket != bucket.BSV {
			continue
		}

		err := info.add(vout, annotation.Paymail)
		if err != nil {
			return info, spverrors.Wrapf(err,
				"failed to process paymail annotation, vout: %d, sender: %s, recipient %s, reference: %s",
				vout, annotation.Paymail.Sender, annotation.Paymail.Receiver, annotation.Paymail.Reference,
			)
		}
	}

	if info.empty() {
		info.skipNotification = true
	}

	return info, nil
}

func (f *txFlow) notifyPaymailExternalRecipient(pmInfo paymailInfo) error {
	if pmInfo.skipNotification {
		f.service.logger.Debug().Str("sender", pmInfo.Sender()).Str("receiver", pmInfo.Receiver()).
			Msg("skipping paymail notification (internal receiver)")
		return nil
	}

	f.service.logger.Info().Str("sender", pmInfo.Sender()).Str("receiver", pmInfo.Receiver()).
		Msg("notifying paymail external recipient")

	err := f.service.paymailNotifier.Notify(
		f.ctx,
		pmInfo.Receiver(),
		&paymail.P2PMetaData{
			Sender: pmInfo.Sender(),
		},
		pmInfo.Reference(),
		f.tx,
	)
	if err != nil {
		return spverrors.Wrapf(err, "failed to notify paymail external recipient")
	}
	return nil
}
