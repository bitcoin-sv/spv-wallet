package main

import (
	"flag"
	"fmt"
)

const (
	domainLocalHost = "localhost:3003"
	adminXPub       = "xpub661MyMwAqRbcFgfmdkPgE2m5UjHXu9dj124DbaGLSjaqVESTWfCD4VuNmEbVPkbYLCkykwVZvmA8Pbf8884TQr1FgdG2nPoHR8aB36YdDQh"
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

}
