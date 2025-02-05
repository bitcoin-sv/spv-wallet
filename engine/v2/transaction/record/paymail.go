package record

import (
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

func processPaymailOutputs(tx *trx.Transaction, userID string, annotations *transaction.Annotations) (*transaction.PaymailAnnotation, error) {
	var paymailAnnotation *transaction.PaymailAnnotation

	for _, annotation := range annotations.Outputs {
		if annotation.Bucket != bucket.BSV {
			continue
		}
		if annotation.Paymail == nil {
			continue
		}

		if paymailAnnotation != nil && *paymailAnnotation != *annotation.Paymail {
			return nil, txerrors.ErrMultiPaymailRecipientsNotSupported
		}

		paymailAnnotation = annotation.Paymail
	}

	// TODO: Get default or check sender paymail

	return paymailAnnotation, nil
}
