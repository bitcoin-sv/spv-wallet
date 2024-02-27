package engine

import (
	"context"
	"errors"
	"fmt"
)

func (c *Client) NewContact(ctx context.Context, fullName, paymail, senderPubKey string, opts ...ModelOps) (*Contact, error) {
	// Check for existing NewRelic transaction
	ctx = c.GetOrStartTxn(ctx, "new_contact")

	contact, err := getContact(ctx, fullName, paymail, senderPubKey, opts...)

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
		)...)...,
	)

	if err != nil {
		return nil, err
	}

	capabilities, err := contact.getContactPaymailCapability(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get contact paymail capability: %w", err)
	}

	pkiURL := capabilities.GetString("pki", "")

	receiverPubKey, err := contact.getPubKeyFromPki(pkiURL)

	contact.PubKey = receiverPubKey

	contact.Status = notConfirmed

	if err = contact.Save(ctx); err != nil {
		return nil, err
	}
	return contact, nil
}
