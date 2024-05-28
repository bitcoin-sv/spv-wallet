package engine

import (
	"context"
	"encoding/hex"
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/spv-wallet/engine/pike"
	"github.com/bitcoin-sv/spv-wallet/engine/script/template"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/libsv/go-bk/bec"
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

func (p *PikePaymentServiceProvider) CreatePikeOutputResponse(
	ctx context.Context,
	alias, domain, senderPubKey string,
	satoshis uint64,
	requestMetadata *server.RequestMetadata,
) (*paymail.PikePaymentOutputsResponse, error) {
	referenceID, err := utils.RandomHex(16)
	if err != nil {
		return nil, err
	}

	outputs, err := pike.GenerateOutputsTemplate(satoshis)
	if err != nil {
		return nil, err
	}

	metadata := createMetadata(requestMetadata, "CreatePikeDestinationResponse")
	opts := WithMetadatas(metadata)

	if err = p.createPikeDestinations(ctx, outputs, alias, domain, senderPubKey, referenceID, opts); err != nil {
		return nil, err
	}

	return &paymail.PikePaymentOutputsResponse{
		Outputs:   convertToPaymailOutputTemplates(outputs),
		Reference: referenceID,
	}, nil
}

func (p *PikePaymentServiceProvider) createPikeDestinations(ctx context.Context, outputsTemplate []*template.OutputTemplate, alias, domain, senderPubKeyHex, reference string, opts ...ModelOps) error {
	pAddress, err := getPaymailAddress(ctx, alias+"@"+domain, p.client.DefaultModelOptions()...)
	if err != nil {
		return err
	}

	receiverPublicKeyHex, err := pAddress.GetPubKey()
	if err != nil {
		return err
	}

	receiverPubKey, senderPubKey, err := getPublicKeys(receiverPublicKeyHex, senderPubKeyHex)
	if err != nil {
		return err
	}

	scripts, err := pike.GenerateLockingScriptsFromTemplates(outputsTemplate, senderPubKey, receiverPubKey, reference)
	if err != nil {
		return err
	}

	return p.saveDestinations(ctx, pAddress, scripts, senderPubKeyHex, opts...)
}

func (p *PikePaymentServiceProvider) saveDestinations(
	ctx context.Context,
	pAddress *PaymailAddress,
	scripts []string,
	senderPubKeyHex string,
	opts ...ModelOps,
) error {
	for index, script := range scripts {
		dst := newDestination(pAddress.XpubID, script, append(p.client.DefaultModelOptions(), opts...)...)
		dst.DerivationMethod = PIKEDerivationMethod
		dst.SenderXpub = senderPubKeyHex
		dst.OutputIndex = uint32(index)

		if err := dst.Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

func getPublicKeys(receiverPubKeyHex, senderPubKeyHex string) (*bec.PublicKey, *bec.PublicKey, error) {
	receiverPubKey, err := getPublicKey(receiverPubKeyHex)
	if err != nil {
		return nil, nil, err
	}

	senderPubKey, err := getPublicKey(senderPubKeyHex)
	if err != nil {
		return nil, nil, err
	}

	return receiverPubKey, senderPubKey, nil
}

func getPublicKey(pubKeyHex string) (*bec.PublicKey, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, err
	}
	return bec.ParsePubKey(pubKeyBytes, bec.S256())
}

func convertToPaymailOutputTemplates(outputTemplates []*template.OutputTemplate) []*paymail.OutputTemplate {
	outputs := make([]*paymail.OutputTemplate, 0)
	for _, output := range outputTemplates {
		outputs = append(outputs, &paymail.OutputTemplate{
			Script:   output.Script,
			Satoshis: output.Satoshis,
		})
	}
	return outputs
}
