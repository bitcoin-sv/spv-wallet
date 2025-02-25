package fixtures

import (
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/custominstructions"
	"strings"

	"github.com/bitcoin-sv/go-paymail"
	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// User is a fixture that is representing a user of the system.
type User struct {
	Paymails []Paymail
	PrivKey  string
}

const (
	// PaymailDomain is the "our" paymail domain in the tests.
	PaymailDomain = "example.com"
	// PaymailDomainExternal is the "their"/external paymail domain in the tests.
	PaymailDomainExternal = "external.example.com"

	// SenderExternalPKI is the PKI of the external Sender used in fixtures
	SenderExternalPKI = "02ed100a85ac774757c967e2a7a8a1c7fdef901795805b494df69d7d02f663d259"
	// RecipientExternalPKI is the PKI of the RecipientExternal used in fixtures
	RecipientExternalPKI = "03bf409b6b2842150142c6b92cb11ba6a06310bdacd0ff2118a9b9da60ed994c2b"
)

var (
	// UserWithMorePaymails is a user with more than one paymail.
	UserWithMorePaymails = User{
		Paymails: []Paymail{
			"tester@" + PaymailDomain,
			"second_pm@" + PaymailDomain,
		},
		PrivKey: "xprv9s21ZrQH143K29ipDWk4vbx6cyyfpbBSj84GrmQPpaKu9Nct6KBhxmSPaGHxoAPisgd3sXKdb2kqKpgLEeAoS54CQGZC8vjoQ6tmJceATxZ",
	}

	// UserWithoutPaymail is a user without any paymail.
	UserWithoutPaymail = User{
		PrivKey: "xprv9s21ZrQH143K4b2JYp37EzEcK55k5wQDnXaH3ooi8oq9yHEj8TCWGuVnJoQvQVyHx3eyF6DyLDiteD6G5CLdKvTcG8QwiEZPyqUcvgmj9aK",
	}

	// Sender is a user that is a sender in the tests.
	Sender = User{
		Paymails: []Paymail{
			"sender@" + PaymailDomain,
		},
		PrivKey: "xprv9s21ZrQH143K2stnKknNEck8NZ9buundyjYCGFGS31bwApaGp7oviHYVY9YAogmgvFC8EdsbsDReydnhDXrRrSXoNoMZczV9t4oPQREAmQ3",
	}

	// RecipientInternal is a user that is a recipient from "our" server in the tests.
	RecipientInternal = User{
		Paymails: []Paymail{
			"recipient@" + PaymailDomain,
		},
		PrivKey: "xprv9s21ZrQH143K3c3jkTBGijY5UsiHUdd3fSzRFD21c7cFduWX4m9nPrcuVrjQ76K234TFWgKF3f97HXggriPipBdhuof6bSvLGE74zCCgJds",
	}

	// RecipientExternal is a user that is a recipient from external server in the tests.
	RecipientExternal = User{
		Paymails: []Paymail{
			"recipient@" + PaymailDomainExternal,
		},
		PrivKey: "xprvA8mj2ZL1w6Nqpi6D2amJLo4Gxy24tW9uv82nQKmamT2rkg5DgjzJZRFnW33e7QJwn65uUWSuN6YQyWrujNjZdVShPRnpNUSRVTru4cxaqfd",
	}

	// SenderExternal is a user that is a sender from external server in the tests.
	SenderExternal = User{
		Paymails: []Paymail{
			"sender@" + PaymailDomainExternal,
		},
		PrivKey: "",
	}

	// ExternalFaucet is a user that is a faucet from external server in the tests.
	ExternalFaucet = User{
		Paymails: []Paymail{
			"faucet@" + PaymailDomainExternal,
		},
		PrivKey: "xprv9s21ZrQH143K3N6qVJQAu4EP51qMcyrKYJLkLgmYXgz58xmVxVLSsbx2DfJUtjcnXK8NdvkHMKfmmg5AJT2nqqRWUrjSHX29qEJwBgBPkJQ",
	}
)

// Paymail wraps a string paymail address to provide additional common methods.
type Paymail string

// String is necessary to fulfill the fmt.Stringer interface.
func (p Paymail) String() string {
	return string(p)
}

// Address returns the address of this paymail.
func (p Paymail) Address() string {
	return string(p)
}

// PublicName returns the public name of this paymail (for testing purposes, it's just the alias in uppercase).
func (p Paymail) PublicName() string {
	return strings.ToUpper(p.Alias())
}

// Domain returns the domain of this paymail.
func (p Paymail) Domain() string {
	_, domain, _ := paymail.SanitizePaymail(p.String())
	return domain
}

// Alias returns the alias of this paymail.
func (p Paymail) Alias() string {
	alias, _, _ := paymail.SanitizePaymail(p.String())
	return alias
}

// DefaultPaymail returns the default paymail of this user.
func (f *User) DefaultPaymail() Paymail {
	if len(f.Paymails) == 0 {
		return ""
	}
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

// XPrivHD returns the xpriv of this user as a HD key.
func (f *User) XPrivHD() *bip32.ExtendedKey {
	xpriv, err := bip32.GenerateHDKeyFromString(f.PrivKey)
	if err != nil {
		panic("Invalid setup of user fixture, cannot restore xpriv: " + err.Error())
	}
	return xpriv
}

// PrivateKey returns the private key of this user.
func (f *User) PrivateKey() *ec.PrivateKey {
	priv, err := bip32.GetPrivateKeyFromHDKey(f.XPrivHD())
	if err != nil {
		panic("Invalid setup of user fixture, cannot restore private key: " + err.Error())
	}
	return priv
}

// PublicKey returns the public key of this user.
func (f *User) PublicKey() *ec.PublicKey {
	return f.PrivateKey().PubKey()
}

// Address returns the address of this user.
func (f *User) Address() *script.Address {
	addr, err := script.NewAddressFromPublicKey(f.PublicKey(), true)
	if err != nil {
		panic("Invalid setup of user fixture, cannot restore address: " + err.Error())
	}
	return addr
}

// ID returns the id of the user.
// Warning: this refers to the new-tx-flow approach where v2 user's id is effectively P2PKH address from the public key.
func (f *User) ID() string {
	return f.Address().AddressString
}

// P2PKHLockingScript returns the locking script of this user.
func (f *User) P2PKHLockingScript(instructions ...bsv.CustomInstruction) *script.Script {
	res, err := custominstructions.NewLockingScriptInterpreter().
		Process(f.PublicKey(), instructions)

	if err != nil {
		panic("Err returned from LockingScriptInterpreter: " + err.Error())
	}

	return res.LockingScript
}

// P2PKHUnlockingScriptTemplate returns the unlocking script template of this user.
func (f *User) P2PKHUnlockingScriptTemplate(instructions ...bsv.CustomInstruction) sdk.UnlockingScriptTemplate {
	res, err := custominstructions.NewInterpreter(&UnlockingTemplateResolver{}).
		Process(f.PrivateKey(), instructions)

	if err != nil {
		panic("Err returned from UnlockingTemplateResolver: " + err.Error())
	}

	return res.Template
}

// AllUsers returns all users fixtures despite it's internal or external user.
func AllUsers() []User {
	return []User{
		UserWithoutPaymail,
		UserWithMorePaymails,
		Sender,
		RecipientInternal,
		RecipientExternal,
		SenderExternal,
		ExternalFaucet,
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
