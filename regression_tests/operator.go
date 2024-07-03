package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

const (
	domainLocalHost    = "localhost:3003"
	adminXPriv         = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
	adminXPub          = "xpub661MyMwAqRbcFgfmdkPgE2m5UjHXu9dj124DbaGLSjaqVESTWfCD4VuNmEbVPkbYLCkykwVZvmA8Pbf8884TQr1FgdG2nPoHR8aB36YdDQh"
	leaderPaymailAlias = "leader"
	minimalBalance     = 100
)

func main() {
	loadConfigFlag := flag.Bool("l", false, "Load configuration from .env.config file")
	flag.Parse()

	var config *Config
	var err error
	user := &User{}

	if *loadConfigFlag {
		config, err = LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}
		user.XPriv = config.ClientOneLeaderXPriv
	} else {
		config = &Config{}
	}

	if !IsSPVWalletRunning(domainLocalHost) {
		fmt.Println("spv-wallet is not running. Run spv-wallet and try again")
		return
	}

	sharedConfig, err := GetSharedConfig(adminXPub)
	if err != nil {
		fmt.Println("Error getting shared config:", err)
		return
	}
	config.ClientOneURL = sharedConfig.PaymailDomains[0]
	config.ClientTwoURL = sharedConfig.PaymailDomains[0]

	if !IsSPVWalletRunning(config.ClientOneURL) {
		fmt.Println("can't connect to spv-wallet at", config.ClientOneURL, "\nEnable tunneling from localhost")
		return
	}

	user, err = handleUserCreation(leaderPaymailAlias, config)
	if err != nil {
		fmt.Println("Error handling user creation:", err)
		return
	}
	UpdateConfigWithUserKeys(config, user)

	if err := handleCoinsTransfer(user, config); err != nil {
		fmt.Println("Error handling transactions:", err)
		return
	}

}

func handleUserCreation(paymailAlias string, config *Config) (*User, error) {
	if config.ClientOneLeaderXPriv != "" {
		if PromptUserAndCheck("Would you like to use user from env? (y/yes or n/no): ") == 1 {
			return UseUserFromEnv(config, paymailAlias)
		}
	}

	user, err := CreateUser(paymailAlias, config)
	if err != nil {
		return handleCreateUserError(err, paymailAlias, config)
	}

	return user, nil
}

func handleCreateUserError(err error, paymailAlias string, config *Config) (*User, error) {
	if err.Error() == "paymail address already exists" {
		return handleExistingPaymail(paymailAlias, config)
	}
	return nil, fmt.Errorf("error creating user: %v", err)
}

func handleExistingPaymail(paymailAlias string, config *Config) (*User, error) {
	if PromptUserAndCheck("Paymail already exists. Would you like to use it (you need to have xpriv)? (y/yes or n/no): ") == 1 {
		return useExistingPaymail(paymailAlias, config)
	}
	if PromptUserAndCheck("Would you like to delete and create new user? (y/yes or n/no):") == 1 {
		return recreateUser(paymailAlias, config)
	}
	return nil, fmt.Errorf("your choices make it impossible to proceed, exiting")
}

func recreateUser(paymailAlias string, config *Config) (*User, error) {
	err := DeleteUser(paymailAlias, config)
	if err != nil {
		return nil, fmt.Errorf("error deleting user: %v", err)
	}
	user, err := CreateUser(paymailAlias, config)
	if err != nil {
		return nil, fmt.Errorf("error creating user after deletion: %v", err)
	}
	return user, nil
}

func useExistingPaymail(paymailAlias string, config *Config) (*User, error) {
	validatedXPriv := GetValidXPriv()
	keys, err := xpriv.FromString(validatedXPriv)
	if err != nil {
		return nil, fmt.Errorf("error parsing xpriv: %v", err)
	}
	return &User{
		XPriv:   keys.XPriv(),
		XPub:    keys.XPub().String(),
		Paymail: PreparePaymail(paymailAlias, config.ClientOneURL),
	}, nil
}

func handleCoinsTransfer(user *User, config *Config) error {
	response := PromptUserAndCheck("Do you have xpriv and master instance URL? (y/yes or n/no): ")
	if response == 0 {
		fmt.Printf("Please send %d Sato for full regression tests:\n%s\n", minimalBalance, user.Paymail)
		isSent := 0
		for isSent < 1 {
			isSent = PromptUserAndCheck("Did you make the transaction? (y/yes or n/no): ")
		}
	} else {
		if err := takeMasterUrlAndXPriv(user); err != nil {
			return fmt.Errorf("error sending coins: %v", err)
		}
	}

	leaderBalance := 0
	for leaderBalance == 0 {
		fmt.Print("Waiting for coins")
		for i := 0; i < 3; i++ {
			fmt.Print(".")
			time.Sleep(1 * time.Second)
		}
		leaderBalance = CheckBalance(config.ClientOneURL, config.ClientOneLeaderXPriv)
		fmt.Println()
	}
	return nil
}

func takeMasterUrlAndXPriv(leaderPaymail *User) error {
	url := GetValidURL()
	xprivMaster := GetValidXPriv()

	err := sendCoinsWithGoClient(url, xprivMaster, leaderPaymail.Paymail)
	if err != nil {
		fmt.Println("Error sending coins:", err)
		return err
	}
	return nil
}

func sendCoinsWithGoClient(instanceUrl string, istanceXPriv string, receiverPaymail string) error {
	client := walletclient.NewWithXPriv(AddPrefixIfNeeded(instanceUrl), istanceXPriv)
	ctx := context.Background()

	balance := CheckBalance(instanceUrl, istanceXPriv)
	if balance < minimalBalance {
		return fmt.Errorf("balance too low: %d", balance)
	}
	recipient := walletclient.Recipients{To: receiverPaymail, Satoshis: uint64(balance - 1)}
	recipients := []*walletclient.Recipients{&recipient}

	_, err := client.SendToRecipients(ctx, recipients, map[string]any{"message": "regression test funds"})
	if err != nil {
		return fmt.Errorf("error sending to recipients: %v", err)
	}
	return nil
}
