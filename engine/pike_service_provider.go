package engine

import (
	"context"
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// PikeContactServiceProvider is an interface for handling the pike contact actions in go-paymail/server
type PikeContactServiceProvider struct {
	client ClientInterface // (pointer) to the Client for accessing SPV Wallet Model methods & etc
}

// PikePaymentServiceProvider is an interface for handling the pike payment actions in go-paymail/server
type PikePaymentServiceProvider struct {
	client ClientInterface // (pointer) to the Client for accessing SPV Wallet Model methods & etc
}

// PIKE CONTACT SERVICE PROVIDER METHODS

func (p *PikeContactServiceProvider) AddContact(
	ctx context.Context,
	requesterPaymailAddress string,
	contact *paymail.PikeContactRequestPayload,
) (err error) {
	if metrics, enabled := p.client.Metrics(); enabled {
		end := metrics.TrackAddContact()
		defer func() {
			success := err == nil
			end(success)
		}()
	}

	reqPaymail, err := getPaymailAddress(ctx, requesterPaymailAddress, p.client.DefaultModelOptions()...)
	if err != nil {
		return
	}
	if reqPaymail == nil {
		err = ErrInvalidRequesterXpub
		return
	}

	_, err = p.client.AddContactRequest(ctx, contact.FullName, contact.Paymail, reqPaymail.XpubID)
	return
}

// PIKE PAYMENT SERVICE PROVIDER METHODS

func (p *PikePaymentServiceProvider) CreatePikeDestinationResponse(
	ctx context.Context,
	alias, domain string,
	satoshis uint64,
	requestMetadata *server.RequestMetadata,
) (*paymail.PikePaymentOutputsResponse, error) {
	referenceID, err := utils.RandomHex(16)
	if err != nil {
		return nil, err
	}

	// Create outputs template
	// outputs := CreatePikeOutputsTemplate(satoshis)

	metadata := createMetadata(requestMetadata, "CreatePikeDestinationResponse")
	opts := WithMetadatas(metadata)

	// Generate and save PIKE destinations
	_, err = p.createPikeDestination(ctx, nil, alias, domain, referenceID, opts)
	if err != nil {
		return nil, err
	}

	return &paymail.PikePaymentOutputsResponse{
		Outputs:   make([]paymail.PikePaymentOutput, 0),
		Reference: referenceID,
	}, nil
}

func (p *PikePaymentServiceProvider) createPikeDestination(ctx context.Context, outputsTemplate []paymail.PikePaymentOutput, alias, domain, reference string, opts ...ModelOps) (*Destination, error) {
	pm, err := getPaymailAddress(ctx, alias+"@"+domain, p.client.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}
	if pm == nil {
		return nil, ErrPaymailNotFound
	}

	hdXpub, err := pm.GetNextXpub(ctx)
	if err != nil {
		return nil, err
	}

	pubKey, err := hdXpub.ECPubKey()
	if err != nil {
		return nil, err
	}

	lockingScript, err := createLockingScript(pubKey)
	if err != nil {
		return nil, err
	}

	// scripts := GenerateLockingScriptsFromTemplates(outputsTemplate, senderPubKey, receiverPubKey, "reference")
	// if len(scripts) == 0 {
	// 	return nil, errors.New("no locking scripts generated")
	// }
	// lockingScript := scripts[0]

	// create a new dst, based on the External xPub child
	// this is not yet possible using the xpub struct. That needs the full xPub, which we don't have.
	dst := newDestination(pm.XpubID, lockingScript, append(opts, New())...)
	dst.Chain = utils.ChainExternal
	dst.Num = pm.ExternalXpubKeyNum
	dst.PaymailExternalDerivationNum = &pm.XpubDerivationSeq

	if err = dst.Save(ctx); err != nil {
		return nil, err
	}
	return dst, nil
}
