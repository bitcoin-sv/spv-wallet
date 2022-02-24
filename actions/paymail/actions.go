package pmail

import (
	"context"
	"time"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux/utils"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bip32"
)

func newPaymailAddress(ctx context.Context, key, address string, opts ...bux.ModelOps) (*PaymailAddress, error) {

	paymailAddress := NewPaymail(address, append(opts, bux.New())...)

	xPub, err := bitcoin.GetHDKeyFromExtendedPublicKey(key)
	if err != nil {
		return nil, err
	}

	var paymailKey *bip32.ExtendedKey
	paymailKey, err = bitcoin.GetHDKeyChild(xPub, utils.ChainExternal)
	if err != nil {
		return nil, err
	}

	paymailAddress.XPubID = utils.Hash(key)
	paymailAddress.ExternalXPubKey = paymailKey.String()

	err = paymailAddress.Save(ctx)
	if err != nil {
		return nil, err
	}

	return paymailAddress, nil
}

func deletePaymailAddress(ctx context.Context, address string, opts ...bux.ModelOps) error {

	paymailAddress, err := GetPaymail(ctx, address, opts...)
	if err != nil {
		return err
	}
	if paymailAddress == nil {
		return ErrMissingPaymail
	}

	var randomString string
	randomString, err = utils.RandomHex(16)
	if err != nil {
		return err
	}

	// We will do a soft delete to make sure we still have the history for this address
	// setting the Domain to a random string solved the problem of the unique index on Alias/Domain
	paymailAddress.Alias = paymailAddress.Alias + "@" + paymailAddress.Domain
	paymailAddress.Domain = randomString
	paymailAddress.DeletedAt.Valid = true
	paymailAddress.DeletedAt.Time = time.Now()

	return paymailAddress.Save(ctx)
}
