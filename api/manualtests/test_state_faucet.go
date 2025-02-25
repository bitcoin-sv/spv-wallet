package manualtests

import (
	"context"
	"strings"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/samber/lo"
)

var FaucetError = StateError.NewSubtype("faucet")

type Faucet struct {
	URL                string
	Xpriv              string
	DefaultTopUpAmount uint64 `mapstructure:"default_topup_amount"`
	state              *State
}

func (f *Faucet) TopUp(satoshis ...uint64) error {
	if len(satoshis) == 0 {
		satoshis = []uint64{f.DefaultTopUpAmount}
	}

	faucet, err := f.UserClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	pub, err := faucet.XPub(ctx)
	if err != nil {
		return FaucetError.Wrap(err, "couldn't check faucet balance")
	}

	logger := Logger().With().Str("faucet", f.URL).Logger()
	logger.Info().Msgf("faucet current balance is %d satoshis", pub.CurrentBalance)

	sum := lo.Sum(satoshis)
	if pub.CurrentBalance < sum+1 {
		return FaucetError.New("faucet balance (%d) is too low to send %d satoshis", pub.CurrentBalance, sum)
	}

	recipients := make([]*commands.Recipients, 0, len(satoshis))
	for i := 0; i < cap(recipients); i++ {
		recipients = append(recipients, &commands.Recipients{
			Satoshis: satoshis[i],
			To:       f.state.CurrentUser().PaymailAddress(),
		})
	}

	tx, err := faucet.SendToRecipients(ctx, &commands.SendToRecipients{
		Metadata: map[string]any{
			"operation":   "v2-paymail-transaction",
			"description": "manual test",
		},
		Recipients: recipients,
	})
	if err != nil {
		return err
	}

	logger.Info().Msgf("Transaction:\n %+v \n", tx)
	return nil
}

func (f *Faucet) UserClient() (*wallet.UserAPI, error) {
	err := f.validate()
	if err != nil {
		return nil, err
	}

	walletClient, err := wallet.NewUserAPIWithXPriv(config.New(config.WithAddr(f.URL)), f.Xpriv)
	if err != nil {
		return nil, FaucetError.Wrap(err, "couldn't create spv-wallet client for user to use it as a faucet")
	}

	return walletClient, nil
}

func (f *Faucet) validate() error {
	if f.URL == "" || f.URL == notConfiguredFaucetURL {
		return FaucetError.New("Configure faucet.url in file://%s", f.state.configFilePath)
	}
	if !strings.HasPrefix(f.Xpriv, "xprv") {
		return FaucetError.New("Configure faucet.xpriv in file://%s", f.state.configFilePath)
	}
	if len(f.Xpriv) != 111 {
		return FaucetError.New("Configured faucet.xpriv too short update file://%s", f.state.configFilePath)
	}

	return nil
}
