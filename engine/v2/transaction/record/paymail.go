package record

import (
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

type paymailInfo struct {
	transaction.PaymailAnnotation
	vouts map[uint32]struct{}
}

func newPaymailInfo(vout uint32, annotation *transaction.PaymailAnnotation) *paymailInfo {
	return &paymailInfo{
		PaymailAnnotation: *annotation,
		vouts:             map[uint32]struct{}{vout: {}},
	}
}

func (pi *paymailInfo) equalsToAnnotation(annotation *transaction.PaymailAnnotation) bool {
	return pi.PaymailAnnotation == *annotation
}

func (pi *paymailInfo) addVOut(vout uint32) {
	pi.vouts[vout] = struct{}{}
}

func (pi *paymailInfo) hasVOut(vout uint32) bool {
	_, ok := pi.vouts[vout]
	return ok
}

func processPaymailOutputs(annotations transaction.Annotations) (*paymailInfo, error) {
	var info *paymailInfo

	for vout, annotation := range annotations.Outputs {
		if annotation.Bucket != bucket.BSV {
			continue
		}
		if annotation.Paymail == nil {
			continue
		}

		if info != nil && !info.equalsToAnnotation(annotation.Paymail) {
			return nil, txerrors.ErrMultiPaymailRecipientsNotSupported
		}

		vout32, err := conv.IntToUint32(vout)
		if err != nil {
			return nil, txerrors.ErrAnnotationIndexConversion.Wrap(err)
		}

		if info == nil {
			info = newPaymailInfo(vout32, annotation.Paymail)
		} else {
			info.addVOut(vout32)
		}
	}

	// TODO: Is "sender" field required?

	return info, nil
}
