package engine

import (
	"context"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// finalizeP2PTransaction will notify the paymail provider about the transaction
func finalizeP2PTransaction(ctx context.Context, client paymail.ClientInterface, p4 *PaymailP4, transaction *Transaction) (*paymail.P2PTransactionPayload, error) {
	if transaction.client != nil {
		transaction.client.Logger().Info().
			Str("txID", transaction.ID).
			Msgf("start %s", p4.Format)
	}

	p2pTransaction, err := buildP2pTx(ctx, p4, transaction)
	if err != nil {
		return nil, err
	}

	response, err := client.SendP2PTransaction(p4.ReceiveEndpoint, p4.Alias, p4.Domain, p2pTransaction)
	if err != nil {
		if transaction.client != nil {
			transaction.client.Logger().Info().
				Str("txID", transaction.ID).
				Msgf("finalizeerror %s, reason: %s", p4.Format, err.Error())
		}
		return nil, spverrors.Wrapf(err, "failed to send transaction via paymail")
	}

	if transaction.client != nil {
		transaction.client.Logger().Info().
			Str("txID", transaction.ID).
			Msgf("successfully finished %s", p4.Format)
	}
	return &response.P2PTransactionPayload, nil
}

func buildP2pTx(ctx context.Context, p4 *PaymailP4, transaction *Transaction) (*paymail.P2PTransaction, error) {
	p2pTransaction := &paymail.P2PTransaction{
		MetaData: &paymail.P2PMetaData{
			Note:   p4.Note,
			Sender: p4.FromPaymail,
		},
		Reference: p4.ReferenceID,
	}

	switch p4.Format {

	case BeefPaymailPayloadFormat:
		beef, err := ToBeef(ctx, transaction, transaction.client)
		if err != nil {
			return nil, err
		}

		p2pTransaction.Beef = beef

	case BasicPaymailPayloadFormat:
		p2pTransaction.Hex = transaction.Hex

	default:
		return nil, spverrors.Newf("%s is unknown format", p4.Format)
	}

	return p2pTransaction, nil
}
