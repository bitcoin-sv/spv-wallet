package engine

import (
	"context"
	"reflect"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/beef"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/go-paymail/spv"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// PaymailDefaultServiceProvider is an interface for overriding the paymail actions in go-paymail/server
//
// This is an example and the default functionality for all the basic Paymail actions
type PaymailDefaultServiceProvider struct {
	client ClientInterface // (pointer) to the Client for accessing SPV Wallet Model methods & etc
}

// createMetadata will create a new metadata seeded from the server information
func createMetadata(serverMetaData *server.RequestMetadata, request string) (metadata Metadata) {
	metadata = make(Metadata)
	metadata["paymail_request"] = request

	if serverMetaData != nil {
		if serverMetaData.UserAgent != "" {
			metadata["user_agent"] = serverMetaData.UserAgent
		}
		if serverMetaData.Note != "" {
			metadata["note"] = serverMetaData.Note
		}
		if serverMetaData.Domain != "" {
			metadata[domainField] = serverMetaData.Domain
		}
		if serverMetaData.IPAddress != "" {
			metadata["ip_address"] = serverMetaData.IPAddress
		}
	}
	return
}

// GetPaymailByAlias will get a paymail address and information by alias
func (p *PaymailDefaultServiceProvider) GetPaymailByAlias(
	ctx context.Context,
	alias, domain string,
	_ *server.RequestMetadata,
) (*paymail.AddressInformation, error) {
	pm, err := getPaymailAddress(ctx, alias+"@"+domain, p.client.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}
	if pm == nil {
		return nil, spverrors.ErrCouldNotFindPaymail
	}

	pk, err := pm.GetPubKey()
	if err != nil {
		return nil, err
	}

	return &paymail.AddressInformation{
		Alias:  pm.Alias,
		Avatar: pm.Avatar,
		Domain: pm.Domain,
		ID:     pm.ID,
		Name:   pm.PublicName,
		PubKey: pk,
	}, nil
}

// CreateAddressResolutionResponse will create the address resolution response
func (p *PaymailDefaultServiceProvider) CreateAddressResolutionResponse(
	ctx context.Context,
	alias, domain string,
	_ bool,
	requestMetadata *server.RequestMetadata,
) (*paymail.ResolutionPayload, error) {
	metadata := createMetadata(requestMetadata, "CreateAddressResolutionResponse")

	dst, err := p.getDestinationForPaymail(ctx, alias, domain, metadata)
	if err != nil {
		return nil, err
	}

	return &paymail.ResolutionPayload{
		Address:   dst.Address,
		Output:    dst.LockingScript,
		Signature: "", // todo: add the signature if senderValidation is enabled
	}, nil
}

// CreateP2PDestinationResponse will create a p2p destination response
func (p *PaymailDefaultServiceProvider) CreateP2PDestinationResponse(
	ctx context.Context,
	alias, domain string,
	satoshis uint64,
	requestMetadata *server.RequestMetadata,
) (*paymail.PaymentDestinationPayload, error) {
	referenceID, err := utils.RandomHex(16)
	if err != nil {
		return nil, spverrors.Wrapf(err, "cannot generate reference id")
	}

	metadata := createMetadata(requestMetadata, "CreateP2PDestinationResponse")
	metadata[ReferenceIDField] = referenceID
	metadata[satoshisField] = satoshis

	// todo: strategy to break apart outputs based on satoshis (return x Outputs)
	dst, err := p.getDestinationForPaymail(ctx, alias, domain, metadata)
	if err != nil {
		return nil, err
	}

	// Append the output(s)
	var outputs []*paymail.PaymentOutput
	outputs = append(outputs, &paymail.PaymentOutput{
		Address:  dst.Address,
		Satoshis: satoshis,
		Script:   dst.LockingScript,
	})

	return &paymail.PaymentDestinationPayload{
		Outputs:   outputs,
		Reference: referenceID,
	}, nil
}

// RecordTransaction will record the transaction
// TODO: rename to HandleReceivedP2pTransaction
func (p *PaymailDefaultServiceProvider) RecordTransaction(ctx context.Context,
	p2pTx *paymail.P2PTransaction, requestMetadata *server.RequestMetadata,
) (*paymail.P2PTransactionPayload, error) {
	// Create the metadata
	metadata := createMetadata(requestMetadata, "HandleReceivedP2pTransaction")
	metadata[p2pMetadataField] = p2pTx.MetaData
	metadata[ReferenceIDField] = p2pTx.Reference

	// Record the transaction
	sdkTx, err := buildSDKTx(p2pTx)
	if err != nil {
		return nil, err
	}

	rts, err := getIncomingTxRecordStrategy(ctx, p.client, sdkTx)
	if err != nil {
		return nil, err
	}
	if err := rts.Validate(); err != nil {
		return nil, err //nolint:wrapcheck // returns our internal errors
	}

	transaction, err := recordTransaction(ctx, p.client, rts, WithMetadatas(metadata))
	if err != nil {
		return nil, err
	}

	if p2pTx.DecodedBeef != nil {
		if reflect.TypeOf(rts) == reflect.TypeOf(&externalIncomingTx{}) {
			go saveBEEFTxInputs(ctx, p.client, p2pTx.DecodedBeef)
		}
	}

	// Return the response from the p2p request
	return &paymail.P2PTransactionPayload{
		Note: p2pTx.MetaData.Note,
		TxID: transaction.ID,
	}, nil
}

// VerifyMerkleRoots will verify the merkle roots by checking them in external header service - Block Headers Service
func (p *PaymailDefaultServiceProvider) VerifyMerkleRoots(
	ctx context.Context,
	merkleRoots []*spv.MerkleRootConfirmationRequestItem,
) (err error) {
	if metrics, enabled := p.client.Metrics(); enabled {
		end := metrics.TrackVerifyMerkleRoots()
		defer func() {
			success := err == nil
			end(success)
		}()
	}

	valid, err := p.client.Chain().VerifyMerkleRoots(ctx, merkleRoots)

	// NOTE: these errors goes to go-paymail and are not logged there, so we need to log them here

	if err != nil {
		p.client.Logger().Error().Err(err).Msg("Error verifying merkle roots")
		return spverrors.ErrPaymailMerkleRootVerificationFailed.Wrap(err)
	}

	if !valid {
		p.client.Logger().Warn().Msg("Not all provided merkle roots were confirmed by BHS")
		return spverrors.ErrPaymailInvalidMerkleRoots
	}

	return
}

// AddContact will add a contact to the paymail address
func (p *PaymailDefaultServiceProvider) AddContact(
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
		err = spverrors.ErrInvalidRequesterXpub
		return
	}

	_, err = p.client.AddContactRequest(ctx, contact.FullName, contact.Paymail, reqPaymail.XpubID)
	return
}

func (p *PaymailDefaultServiceProvider) getDestinationForPaymail(ctx context.Context, alias, domain string, metadata Metadata) (*Destination, error) {
	pm, err := getPaymailAddress(ctx, alias+"@"+domain, p.client.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}
	if pm == nil {
		return nil, spverrors.ErrCouldNotFindPaymail
	}

	dst, err := createDestination(
		ctx, pm, append(p.client.DefaultModelOptions(), WithMetadatas(metadata))...,
	)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func createDestination(ctx context.Context, pm *PaymailAddress, opts ...ModelOps) (*Destination, error) {
	hdXpub, err := pm.GetNextXpub(ctx)
	if err != nil {
		return nil, err
	}

	pubKey, err := hdXpub.ECPubKey()
	if err != nil {
		return nil, spverrors.Wrapf(err, "fail to get next public key for destination")
	}

	lockingScript, err := createLockingScript(pubKey)
	if err != nil {
		return nil, err
	}

	// create a new dst, based on the External xPub child
	// this is not yet possible using the xpub struct. That needs the full xPub, which we don't have.
	dst := newDestination(pm.XpubID, lockingScript, append(opts, New())...)
	dst.Chain = utils.ChainExternal
	dst.Num = pm.ExternalXpubKeyNum
	dst.PaymailExternalDerivationNum = &pm.XpubDerivationSeq
	dst.DerivationMethod = BIP32DerivationMethod

	if err = dst.Save(ctx); err != nil {
		return nil, err
	}

	return dst, nil
}

func createLockingScript(ecPubKey *ec.PublicKey) (lockingScript string, err error) {
	bsvAddress, err := script.NewAddressFromPublicKey(ecPubKey, true)
	if err != nil {
		return
	}

	ls, err := p2pkh.Lock(bsvAddress)
	lockingScript = ls.String()
	return
}

func buildSDKTx(p2pTx *paymail.P2PTransaction) (*trx.Transaction, error) {
	var err error
	var tx *trx.Transaction
	if p2pTx.Beef != "" {
		tx, err = trx.NewTransactionFromBEEFHex(p2pTx.Beef)
		if err != nil {
			return nil, spverrors.Wrapf(err, "unable to create transaction from BEEF")
		}

		return tx, nil
	}
	tx, err = trx.NewTransactionFromHex(p2pTx.Hex)
	if err != nil {
		return nil, spverrors.Wrapf(err, "unable to create transaction from hex")
	}

	return tx, nil
}

func saveBEEFTxInputs(ctx context.Context, c ClientInterface, dBeef *beef.DecodedBEEF) {
	inputsToAdd, err := getInputsWhichAreNotInDb(c, dBeef)
	if err != nil {
		c.Logger().Error().Msgf("error in saveBEEFTxInputs: %v", err)
	}

	for _, input := range inputsToAdd {
		var bump *BUMP
		if input.BumpIndex != nil { // mined
			bumpIndex := int(*input.BumpIndex) //nolint:gosec // we trust the input
			bump, err = getBump(bumpIndex, dBeef.BUMPs)
			if err != nil {
				c.Logger().Error().Msgf("error in saveBEEFTxInputs: %v for beef: %v", err, dBeef)
			}

		}

		err = saveBeefTransactionInput(ctx, c, input, bump)
		if err != nil {
			c.Logger().Error().Msgf("error in saveBEEFTxInputs: %v for beef: %v", err, dBeef)
		}
	}
}

func getInputsWhichAreNotInDb(c ClientInterface, dBeef *beef.DecodedBEEF) ([]*beef.TxData, error) {
	txIDs := make([]string, 0, len(dBeef.Transactions))
	for _, tx := range dBeef.Transactions {
		txIDs = append(txIDs, tx.GetTxID())
	}
	dbTxs, err := c.GetTransactionsByIDs(context.Background(), txIDs)
	if err != nil {
		return nil, spverrors.Wrapf(err, "error during getting txs from db")
	}

	txs := make([]*beef.TxData, 0)

	if len(dbTxs) == len(txIDs) {
		return txs, nil
	}

	for _, input := range dBeef.Transactions {
		found := false
		for _, dbTx := range dbTxs {
			if dbTx.GetID() == input.GetTxID() {
				found = true
				break
			}
		}
		if !found {
			txs = append(txs, input)
		}
	}

	return txs, nil
}

func getBump(bumpIndex int, bumps beef.BUMPs) (*BUMP, error) {
	if bumpIndex > len(bumps) {
		return nil, spverrors.Newf("error in getBump: bump index exceeds bumps length")
	}

	bump := bumps[bumpIndex]
	paths := make([][]BUMPLeaf, 0)

	for _, path := range bump.Path {
		pathLeaves := make([]BUMPLeaf, 0)
		for _, leaf := range path {
			l := BUMPLeaf{
				Offset:    leaf.Offset,
				Hash:      leaf.Hash,
				TxID:      leaf.TxId,
				Duplicate: leaf.Duplicate,
			}
			pathLeaves = append(pathLeaves, l)
		}
		paths = append(paths, pathLeaves)
	}

	return &BUMP{
		BlockHeight: bump.BlockHeight,
		Path:        paths,
	}, nil
}

func saveBeefTransactionInput(ctx context.Context, c ClientInterface, input *beef.TxData, bump *BUMP) error {
	newOpts := c.DefaultModelOptions(New())
	inputTx, _ := txFromHex(input.Transaction.String(), newOpts...) // we can ignore error here

	if bump != nil {
		inputTx.BUMP = *bump
		inputTx.TxStatus = TxStatusMined
	}

	err := inputTx.Save(ctx)
	return spverrors.Wrapf(err, "error in saveBeefTransactionInput during saving tx")
}
