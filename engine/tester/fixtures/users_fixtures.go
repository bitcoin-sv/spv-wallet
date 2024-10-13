package fixtures

import (
	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// User is a fixture that is representing a user of the system.
type User struct {
	Paymails []string
	PrivKey  string
}

const (
	// PaymailDomain is the "our" paymail domain in the tests.
	PaymailDomain = "example.com"
	// PaymailDomainExternal is the "their"/external paymail domain in the tests.
	PaymailDomainExternal = "external.example.com"
)

var (
	// UserWithMorePaymails is a user with more than one paymail.
	UserWithMorePaymails = User{
		Paymails: []string{
			"tester@" + PaymailDomain,
			"secondPm@" + PaymailDomain,
		},
		PrivKey: "xprv9s21ZrQH143K29ipDWk4vbx6cyyfpbBSj84GrmQPpaKu9Nct6KBhxmSPaGHxoAPisgd3sXKdb2kqKpgLEeAoS54CQGZC8vjoQ6tmJceATxZ",
	}

	// UserWithoutPaymail is a user without any paymail.
	UserWithoutPaymail = User{
		PrivKey: "xprv9s21ZrQH143K4b2JYp37EzEcK55k5wQDnXaH3ooi8oq9yHEj8TCWGuVnJoQvQVyHx3eyF6DyLDiteD6G5CLdKvTcG8QwiEZPyqUcvgmj9aK",
	}

	// Sender is a user that is a sender in the tests.
	Sender = User{
		Paymails: []string{
			"sender@" + PaymailDomain,
		},
		PrivKey: "xprv9s21ZrQH143K2stnKknNEck8NZ9buundyjYCGFGS31bwApaGp7oviHYVY9YAogmgvFC8EdsbsDReydnhDXrRrSXoNoMZczV9t4oPQREAmQ3",
	}

	// RecipientInternal is a user that is a recipient from "our" server in the tests.
	RecipientInternal = User{
		Paymails: []string{
			"recipient@" + PaymailDomain,
		},
		PrivKey: "xprv9s21ZrQH143K3c3jkTBGijY5UsiHUdd3fSzRFD21c7cFduWX4m9nPrcuVrjQ76K234TFWgKF3f97HXggriPipBdhuof6bSvLGE74zCCgJds",
	}

	// RecipientExternal is a user that is a recipient from external server in the tests.
	RecipientExternal = User{
		Paymails: []string{
			"recipient@" + PaymailDomainExternal,
		},
		PrivKey: "",
	}
)

// DefaultPaymail returns the default paymail of this user.
func (f *User) DefaultPaymail() string {
	return f.Paymails[0]
}

// XPriv returns the xpriv of this user.
func (f *User) XPriv() string {
	return f.PrivKey
}

// XPub returns the xpub of this user.
// We're calculating it to avoid mistakes in setting up the fixtures.
func (f *User) XPub() string {
	if f.PrivKey == "" {
		return ""
	}
	key, err := bip32.NewKeyFromString(f.PrivKey)
	if err != nil {
		panic("Invalid setup of user fixture, cannot restore xpriv: " + err.Error())
	}
	pubkey, err := key.Neuter()
	if err != nil {
		panic("Invalid setup of user fixture, cannot calculate xpub: " + err.Error())
	}
	return pubkey.String()
}

// XPubID returns the xpub id of this user.
// We're calculating it to avoid mistakes in setting up the fixtures.
func (f *User) XPubID() string {
	xpub := f.XPub()
	if xpub == "" {
		return ""
	}
	return utils.Hash(xpub)
}

// AllUsers returns all users fixtures despite it's internal or external user.
func AllUsers() []User {
	return []User{
		UserWithoutPaymail,
		UserWithMorePaymails,
		Sender,
		RecipientInternal,
		RecipientExternal,
	}
}

// InternalUsers returns all users fixtures representing spv-wallet users.
func InternalUsers() []User {
	return []User{
		UserWithoutPaymail,
		UserWithMorePaymails,
		Sender,
		RecipientInternal,
	}
}
