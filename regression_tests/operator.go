package main

import (
	"flag"
	"fmt"
)

const (

	domainLocalHost = "localhost:3003"
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

}
