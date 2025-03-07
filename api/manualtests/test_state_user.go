package manualtests

import (
	"fmt"
	"time"

	compat "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/samber/lo"
)

var UserDeleted = StateError.NewSubtype("userDeleted")

type AdditionalAlias struct {
	Alias string `mapstructure:"alias" yaml:"alias"`
}

func (a *AdditionalAlias) String() string {
	return a.Alias
}

func (a *AdditionalAlias) PublicName() string {
	return lo.Capitalize(a.Alias)
}

type User struct {
	Alias             string            `mapstructure:"alias" yaml:"alias"`
	Domain            string            `mapstructure:"domain" yaml:"domain"`
	Xpriv             string            `mapstructure:"xpriv" yaml:"xpriv"`
	Xpub              string            `mapstructure:"xpub" yaml:"xpub"`
	PrivateKey        string            `mapstructure:"private_key" yaml:"private_key"`
	PublicKey         string            `mapstructure:"public_key" yaml:"public_key"`
	ID                string            `mapstructure:"id" yaml:"id"`
	AdditionalAliases []AdditionalAlias `mapstructure:"additional_aliases" yaml:"additional_aliases"`
	DataOutpoints     []string          `mapstructure:"data_outpoints" yaml:"data_outpoints"`
	Tags              []string          `mapstructure:"tags" yaml:"tags"`

	// Note holds your notes about the user
	// This field prevent removing your custom notes when updating a state in file.
	// Notice the underscore in _note, this ensures that the note will be top field in the yaml file.
	Note string `mapstructure:"_note" yaml:"_note"`

	state   *State
	xpriv   *compat.ExtendedKey
	xpub    *compat.ExtendedKey
	privKey *ec.PrivateKey
	pubkey  *ec.PublicKey
}

func (u *User) IsEmpty() bool {
	return u.Xpriv == "" || u.Alias == "" || u.Domain == ""
}

func (u *User) new(xpriv string, xpub string) error {
	now := time.Now()

	user := &User{
		Xpriv:             xpriv,
		Xpub:              xpub,
		AdditionalAliases: make([]AdditionalAlias, 0),
		// format now to yymmddhhmmss
		Alias:  "test" + now.Format("060102150405"),
		Domain: u.state.Domain,
	}

	err := user.init()
	if err != nil {
		return err
	}
	*u = *user

	return nil
}

func (u *User) GetPrivateKey() *ec.PrivateKey {
	if u.privKey == nil {
		panic("user wasn't initialized")
	}
	return u.privKey
}

func (u *User) PaymailAddress() string {
	return u.Alias + "@" + u.Domain
}

func (u *User) ShouldGetPaymailAddress() (*string, error) {
	if u.IsEmpty() {
		return nil, StateError.New("there is no current user, before using this method create a user as admin first.")
	}

	return lo.ToPtr(u.PaymailAddress()), nil
}

func (u *User) ShouldGetAdditionalPaymailAddress() (*string, error) {
	if u.IsEmpty() {
		return nil, StateError.New("there is no current user, before using this method create a user as admin first.")
	}

	if len(u.AdditionalAliases) == 0 {
		return nil, StateError.New("there is no additional paymail address, before using this method add an additional paymail address first with admin API.")
	}

	additionalAddress := u.AdditionalAliases[0].String() + "@" + u.Domain
	return lo.ToPtr(additionalAddress), nil
}

func (u *User) PublicName() string {
	return lo.Capitalize(u.Alias)
}

func (u *User) Address() string {
	return u.ID
}

func (u *User) AvatarURL() string {
	return "https://rd.centraltest.com/tests/AVATARTestLogos/TG_Avatar_Logo.png"
}

func (u *User) address() *script.Address {
	pubKey, err := ec.PublicKeyFromString(u.PublicKey)
	if err != nil {
		err = StateError.Wrap(err, "could not get public key from user")
		panic(err)
	}

	addr, err := script.NewAddressFromPublicKey(pubKey, true)
	if err != nil {
		err = StateError.Wrap(err, "could not get address from public key")
		panic(err)
	}
	return addr
}

func (u *User) CreateAdditionalAlias() AdditionalAlias {
	now := time.Now()

	prefix := u.Alias

	// format now to yymmddhhmmss
	additional := AdditionalAlias{
		Alias: prefix + "+" + now.Format("060102150405"),
	}

	u.AdditionalAliases = append(u.AdditionalAliases, additional)

	return additional
}

func (u *User) init() error {
	if u.Xpriv == "" {
		return nil
	}

	xpriv, err := compat.NewKeyFromString(u.Xpriv)
	if err != nil {
		return StateError.Wrap(err, "could not get xpriv from string")
	}

	xpub, err := xpriv.Neuter()
	if err != nil {
		return StateError.Wrap(err, "could not get xpub from xpriv")
	}

	privKey, err := xpriv.ECPrivKey()
	if err != nil {
		return StateError.Wrap(err, "could not get private key from xpriv")
	}
	pubKey, err := xpriv.ECPubKey()
	if err != nil {
		return StateError.Wrap(err, "could not get public key from xpriv")
	}

	u.Xpub = xpub.String()
	u.PrivateKey = privKey.Wif()
	u.PublicKey = pubKey.ToDERHex()

	u.xpriv = xpriv
	u.xpub = xpub
	u.privKey = privKey
	u.pubkey = pubKey

	u.ID = u.address().AddressString

	return nil
}

func (u *User) AddDataOutpoint(txID string, vout int) {
	u.DataOutpoints = append(u.DataOutpoints, fmt.Sprintf("%s-%d", txID, vout))
}

func (u *User) MarkAsDeleted() error {
	if lo.Contains(u.DataOutpoints, "deleted") {
		return UserDeleted.New("user is already marked as deleted")
	}
	u.Tags = append(u.Tags, "deleted")
	return nil
}

func (u *User) RemoveTag(s string) {
	u.Tags = lo.Filter(u.Tags, func(tag string, _ int) bool {
		return tag != s
	})
}
