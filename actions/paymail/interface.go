package pmail

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/config"
	"github.com/BuxOrg/bux/utils"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/tonicpow/go-paymail"
	"github.com/tonicpow/go-paymail/server"
)

const (
	paymailMetadataField    = "metadata"
	paymailRequestField     = "paymail_request"
	paymailP2PMetadataField = "p2p_tx_metadata"
	defaultGetTimeout       = 10 * time.Second
)

// PaymailInterface is an interface for overriding the paymail functions
type PaymailInterface struct {
	client    bux.ClientInterface
	appConfig *config.AppConfig
}

// addServerMetadata adds server data to the existing metadata
func addServerMetadata(serverMetaData *server.RequestMetadata, metadata bux.Metadata) {
	if serverMetaData != nil {
		if serverMetaData.UserAgent != "" {
			metadata["user-agent"] = serverMetaData.UserAgent
		}
		if serverMetaData.Note != "" {
			metadata["note"] = serverMetaData.Note
		}
	}
}

// NewServiceProvider makes a new service provider interface
func NewServiceProvider(buxEngine bux.ClientInterface, config *config.AppConfig) server.PaymailServiceProvider {
	return &PaymailInterface{
		client:    buxEngine,
		appConfig: config,
	}
}

// GetPaymailByAlias will get a paymail address and information by alias
func (p *PaymailInterface) GetPaymailByAlias(ctx context.Context, alias, domain string,
	requestMetadata *server.RequestMetadata) (*paymail.AddressInformation, error) {

	metadata := make(bux.Metadata)
	metadata[paymailRequestField] = "GetPaymailByAlias"
	addServerMetadata(requestMetadata, metadata)
	paymailAddress, pubKey, destination, err := p.getPaymailInformation(ctx, alias, domain, &metadata)
	if err != nil {
		return nil, err
	}

	return &paymail.AddressInformation{
		Alias:       paymailAddress.Alias,
		Avatar:      paymailAddress.Avatar,
		Domain:      paymailAddress.Domain,
		ID:          paymailAddress.ID,
		LastAddress: destination.Address,
		Name:        paymailAddress.Username,
		PrivateKey:  "",
		PubKey:      pubKey,
	}, nil
}

// CreateAddressResolutionResponse will create the address resolution response
func (p *PaymailInterface) CreateAddressResolutionResponse(ctx context.Context, alias, domain string,
	_ bool, requestMetadata *server.RequestMetadata) (*paymail.ResolutionPayload, error) {

	metadata := make(bux.Metadata)
	metadata[paymailRequestField] = "CreateAddressResolutionResponse"
	addServerMetadata(requestMetadata, metadata)
	_, _, destination, err := p.getPaymailInformation(ctx, alias, domain, &metadata)
	if err != nil {
		return nil, err
	}

	return &paymail.ResolutionPayload{
		Address:   destination.Address,
		Output:    destination.LockingScript,
		Signature: "", // todo: add the signature if senderValidation is enabled
	}, nil
}

// CreateP2PDestinationResponse will create a p2p destination response
func (p *PaymailInterface) CreateP2PDestinationResponse(ctx context.Context, alias, domain string,
	satoshis uint64, requestMetadata *server.RequestMetadata) (*paymail.PaymentDestinationPayload, error) {

	referenceID, err := utils.RandomHex(12)
	if err != nil {
		return nil, err
	}

	var outputs []*paymail.PaymentOutput
	outputs, err = p.getOutputDestinations(ctx, alias, domain, satoshis, referenceID, requestMetadata)
	if err != nil {
		return nil, err
	}

	return &paymail.PaymentDestinationPayload{
		Outputs:   outputs,
		Reference: referenceID,
	}, nil
}

func (p *PaymailInterface) getOutputDestinations(ctx context.Context, alias, domain string,
	satoshis uint64, referenceID string, requestMetadata *server.RequestMetadata) ([]*paymail.PaymentOutput, error) {

	var outputs []*paymail.PaymentOutput

	// todo split destinations up into multiple
	metadata := make(bux.Metadata)
	metadata[paymailRequestField] = "CreateP2PDestinationResponse"
	metadata[bux.ReferenceIDField] = referenceID
	metadata["satoshis"] = satoshis
	addServerMetadata(requestMetadata, metadata)

	// this could be improved by creating multiple destinations, for privacy reasons
	_, _, destination, err := p.getPaymailInformation(ctx, alias, domain, &metadata)
	if err != nil {
		return nil, err
	}

	outputs = append(outputs, &paymail.PaymentOutput{
		Address:  destination.Address,
		Satoshis: satoshis,
		Script:   destination.LockingScript,
	})

	return outputs, nil
}

// RecordTransaction will record the transaction
func (p *PaymailInterface) RecordTransaction(ctx context.Context,
	p2pTx *paymail.P2PTransaction, metaData *server.RequestMetadata) (*paymail.P2PTransactionPayload, error) {

	if p.appConfig.GDPRCompliance {
		// remove PII that cannot be in here for GDPR reasons
		metaData.IPAddress = ""
	}

	var opts []bux.ModelOps
	opts = append(opts,
		bux.WithMetadatas(map[string]interface{}{
			paymailRequestField:     "RecordTransaction",
			paymailMetadataField:    metaData,
			paymailP2PMetadataField: p2pTx.MetaData,
			bux.ReferenceIDField:    p2pTx.Reference,
		}),
	)

	transaction, err := p.client.RecordTransaction(
		ctx, "", p2pTx.Hex, "", opts...,
	)
	if err != nil {
		return nil, err
	}

	return &paymail.P2PTransactionPayload{
		Note: p2pTx.MetaData.Note,
		TxID: transaction.ID,
	}, nil
}

// getPaymailInformation will get the paymail information
func (p *PaymailInterface) getPaymailInformation(ctx context.Context, alias,
	domain string, metadata *bux.Metadata) (*bux.PaymailAddress, string, *bux.Destination, error) {

	// todo xPub locking?

	paymailAddress, err := p.getPaymailAddress(ctx, alias, domain)
	if err != nil {
		return nil, "", nil, err
	}

	var xPub *bux.Xpub
	if xPub, err = p.client.GetXpubByID(
		ctx, paymailAddress.XpubID,
	); err != nil {
		return nil, "", nil, err
	}

	pubKey, address, lockingScript, keyErr := p.getPaymailKeys(
		paymailAddress.ExternalXpubKey,
		xPub.NextExternalNum,
	)
	if keyErr != nil {
		return nil, "", nil, keyErr
	}

	// create a new destination, based on the External xPub child
	// this is not yet possible within this library, it needs the full xPub
	destination := &bux.Destination{
		Model: *bux.NewBaseModel(
			bux.ModelDestination,
			append(p.client.DefaultModelOptions(), bux.WithMetadatas(*metadata))...,
		),
		Address:       address,
		Chain:         utils.ChainExternal,
		ID:            utils.Hash(lockingScript),
		LockingScript: lockingScript,
		Num:           xPub.NextExternalNum,
		Type:          utils.ScriptTypePubKeyHash,
		XpubID:        paymailAddress.XpubID,
	}
	if len(*metadata) > 0 {
		destination.Metadata = *metadata
	}

	if err = destination.Save(ctx); err != nil {
		return nil, "", nil, err
	}

	xPub.NextExternalNum++

	if err = xPub.Save(ctx); err != nil {
		return nil, "", nil, err
	}

	return paymailAddress, pubKey, destination, nil
}

// getPaymailAddress will get a paymail address
func (p *PaymailInterface) getPaymailAddress(ctx context.Context, alias, domain string) (*bux.PaymailAddress, error) {
	paymailAddress := &bux.PaymailAddress{
		Model: *bux.NewBaseModel(bux.ModelPaymail, p.client.DefaultModelOptions()...),
	}
	conditions := map[string]interface{}{
		"alias":  alias,
		"domain": domain,
	}
	if err := bux.Get(
		ctx, paymailAddress, conditions, false, defaultGetTimeout,
	); err != nil {
		return nil, err
	}
	if paymailAddress.ExternalXpubKey == "" {
		// could not find paymail record
		return nil, bux.ErrMissingXpub
	}

	return paymailAddress, nil
}

// getPaymailKeys will get all the paymail keys
func (p *PaymailInterface) getPaymailKeys(rawXPubKey string, num uint32) (string, string, string, error) {
	hdKey, err := utils.ValidateXPub(rawXPubKey)
	if err != nil {
		return "", "", "", err
	}

	var derivedKey *bip32.ExtendedKey
	if derivedKey, err = bitcoin.GetHDKeyChild(hdKey, num); err != nil {
		return "", "", "", err
	}

	var nextKey *bec.PublicKey
	if nextKey, err = derivedKey.ECPubKey(); err != nil {
		return "", "", "", err
	}
	pubKey := hex.EncodeToString(nextKey.SerialiseCompressed())

	var bsvAddress *bscript.Address
	if bsvAddress, err = bitcoin.GetAddressFromPubKey(
		nextKey, true,
	); err != nil {
		return "", "", "", err
	}
	address := bsvAddress.AddressString

	var lockingScript string
	if lockingScript, err = bitcoin.ScriptFromAddress(address); err != nil {
		return "", "", "", err
	}

	return pubKey, address, lockingScript, nil
}
