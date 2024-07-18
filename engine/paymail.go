package engine

import (
	"context"
	"errors"
	"strings"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/mrz1836/go-cachestore"
)

// getCapabilities is a utility function to retrieve capabilities for a Paymail provider
func getCapabilities(ctx context.Context, cs cachestore.ClientInterface, client paymail.ClientInterface,
	domain string,
) (*paymail.CapabilitiesPayload, error) {
	// Attempt to get from cachestore
	// todo: allow user to configure the time that they want to cache the capabilities (if they want to cache or not)
	capabilities := new(paymail.CapabilitiesPayload)
	if err := cs.GetModel(
		ctx, cacheKeyCapabilities+domain, capabilities,
	); err != nil && !errors.Is(err, cachestore.ErrKeyNotFound) {
		return nil, spverrors.Wrapf(err, "failed to get capabilities from cachestore")
	} else if len(capabilities.Capabilities) > 0 {
		return capabilities, nil
	}

	// Get SRV record (domain can be different!)
	var response *paymail.CapabilitiesResponse
	srv, err := client.GetSRVRecord(
		paymail.DefaultServiceName, paymail.DefaultProtocol, domain,
	)
	if err != nil {
		// Error returned was a real error
		if !strings.Contains(err.Error(), "zero SRV records found") { // This error is from no SRV record being found
			return nil, err //nolint:wrapcheck // we have handler for paymail errors
		}

		// Try to get capabilities without the SRV record
		// 'Should no record be returned, a paymail client should assume a host of <domain>.<tld> and a port of 443.'
		// http://bsvalias.org/02-01-host-discovery.html

		// Get the capabilities via target
		if response, err = client.GetCapabilities(
			domain, paymail.DefaultPort,
		); err != nil {
			return nil, err //nolint:wrapcheck // we have handler for paymail errors
		}
	} else {
		// Get the capabilities via SRV record
		if response, err = client.GetCapabilities(
			srv.Target, int(srv.Port),
		); err != nil {
			return nil, err //nolint:wrapcheck // we have handler for paymail errors
		}
	}

	// Save to cachestore
	if cs != nil && !cs.Engine().IsEmpty() {
		_ = cs.SetModel(
			context.Background(), cacheKeyCapabilities+domain,
			&response.CapabilitiesPayload, cacheTTLCapabilities,
		)
	}

	return &response.CapabilitiesPayload, nil
}

// hasP2P will return the P2P urls and true if they are both found
func hasP2P(capabilities *paymail.CapabilitiesPayload) (success bool, p2pDestinationURL, p2pSubmitTxURL string, format PaymailPayloadFormat) {
	p2pDestinationURL = capabilities.GetString(paymail.BRFCP2PPaymentDestination, "")
	p2pSubmitTxURL = capabilities.GetString(paymail.BRFCP2PTransactions, "")
	p2pBeefSubmitTxURL := capabilities.GetString(paymail.BRFCBeefTransaction, "")

	if len(p2pBeefSubmitTxURL) > 0 {
		p2pSubmitTxURL = p2pBeefSubmitTxURL
		format = BeefPaymailPayloadFormat
	}

	if len(p2pSubmitTxURL) > 0 && len(p2pDestinationURL) > 0 {
		success = true
	}
	return
}

// startP2PTransaction will start the P2P transaction, returning the reference ID and outputs
func startP2PTransaction(client paymail.ClientInterface,
	alias, domain, p2pDestinationURL string, satoshis uint64,
) (*paymail.PaymentDestinationPayload, error) {
	// Start the P2P transaction request
	response, err := client.GetP2PPaymentDestination(
		p2pDestinationURL,
		alias, domain,
		&paymail.PaymentRequest{Satoshis: satoshis},
	)
	if err != nil {
		return nil, err //nolint:wrapcheck // we have handler for paymail errors
	}

	return &response.PaymentDestinationPayload, nil
}

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
