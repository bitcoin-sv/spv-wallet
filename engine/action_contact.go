package engine

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitcoin-sv/go-paymail"
	"github.com/mrz1836/go-cachestore"
	"github.com/mrz1836/go-datastore"
)

func (c *Client) NewContact(ctx context.Context, fullName, paymail, senderPubKey string, opts ...ModelOps) (*Contact, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "new_contact")

	contact, err := getContact(ctx, fullName, paymail, senderPubKey, c.DefaultModelOptions()...)

	if contact != nil {
		return nil, errors.New("contact already exists")
	}
	if err != nil {
		return nil, err
	}

	contact, err = newContact(
		fullName,
		paymail,
		senderPubKey,
		append(opts, c.DefaultModelOptions(
			New(),
			WithXPub(senderPubKey),
		)...)...,
	)

	if err != nil {
		return nil, err
	}

	capabilities, err := c.GetPaymailCapability(ctx, contact.Paymail)

	if err != nil {
		return nil, fmt.Errorf("failed to get contact paymail capability: %w", err)
	}

	pkiURL := capabilities.GetString("pki", "")

	receiverPubKey, err := c.GetPubKeyFromPki(pkiURL, contact.Paymail)

	contact.PubKey = receiverPubKey

	contact.Status = notConfirmed

	if err = contact.Save(ctx); err != nil {
		return nil, err
	}
	return contact, nil
}

func (c *Client) GetPubKeyFromPki(pkiUrl, paymailAddress string) (string, error) {
	if pkiUrl == "" {
		return "", errors.New("pkiUrl should not be empty")
	}
	alias, domain, _ := paymail.SanitizePaymail(paymailAddress)
	pc := c.PaymailClient()

	pkiResponse, err := pc.GetPKI(pkiUrl, alias, domain)

	if err != nil {
		return "", fmt.Errorf("error getting public key from PKI: %w", err)
	}
	return pkiResponse.PubKey, nil
}

func (c *Client) GetPaymailCapability(ctx context.Context, paymailAddress string) (*paymail.CapabilitiesPayload, error) {
	address := newPaymail(paymailAddress)

	cs := c.Cachestore()
	pc := c.PaymailClient()

	capabilities, err := getCapabilities(ctx, cs, pc, address.Domain)

	if err != nil {
		if errors.Is(err, cachestore.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return capabilities, nil
}

func (c *Client) GetContacts(ctx context.Context, metadata *Metadata, conditions *map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps) ([]*Contact, error) {

	ctx = c.GetOrStartTxn(ctx, "get_contacts")

	contacts, err := getContacts(ctx, metadata, conditions, queryParams, c.DefaultModelOptions(opts...)...)

	if err != nil {
		return nil, err
	}

	return contacts, nil
}
