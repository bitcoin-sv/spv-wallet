package main

import (
	"flag"
	"fmt"
)

const (
	ClientOneURLEnvVar         = "CLIENT_ONE_URL"
	ClientTwoURLEnvVar         = "CLIENT_TWO_URL"
	ClientOneLeaderXPrivEnvVar = "CLIENT_ONE_LEADER_XPRIV"
	ClientTwoLeaderXPrivEnvVar = "CLIENT_TWO_LEADER_XPRIV"
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
}
