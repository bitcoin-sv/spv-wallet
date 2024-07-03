package main

import (
	"flag"
	"fmt"
)

const (
	domainLocalHost    = "localhost:3003"
	adminXPub          = "xpub661MyMwAqRbcFgfmdkPgE2m5UjHXu9dj124DbaGLSjaqVESTWfCD4VuNmEbVPkbYLCkykwVZvmA8Pbf8884TQr1FgdG2nPoHR8aB36YdDQh"
	leaderPaymailAlias = "leader"
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

}

func handleUserCreation(paymailAlias string, config *Config) (*User, error) {
	if config.ClientOneLeaderXPriv != "" {
		if PromptUserAndCheck("Would you like to use user from env? (y/yes or n/no): ") == 1 {
			return useUserFromEnv(config, paymailAlias)
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
		return nil, nil //TODO: Add logic to handle existing user
	}
	return nil, fmt.Errorf("error creating user: %v", err)
}
