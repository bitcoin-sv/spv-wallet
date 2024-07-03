package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type User struct {
	XPriv   string `json:"xpriv"`
	XPub    string `json:"xpub"`
	Paymail string `json:"paymail"`
}

type Config struct {
	ClientOneURL         string
	ClientTwoURL         string
	ClientOneLeaderXPriv string
	ClientTwoLeaderXPriv string
}

func SaveConfig(config *Config) error {
	envMap := map[string]string{
		ClientOneURLEnvVar:         config.ClientOneURL,
		ClientTwoURLEnvVar:         config.ClientTwoURL,
		ClientOneLeaderXPrivEnvVar: config.ClientOneLeaderXPriv,
		ClientTwoLeaderXPrivEnvVar: config.ClientTwoLeaderXPriv,
	}

	err := godotenv.Write(envMap, ".env")
	if err != nil {
		return fmt.Errorf("error saving .env file: %w", err)
	}

	return nil
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return &Config{
		ClientOneURL:         os.Getenv(ClientOneURLEnvVar),
		ClientTwoURL:         os.Getenv(ClientTwoURLEnvVar),
		ClientOneLeaderXPriv: os.Getenv(ClientOneLeaderXPrivEnvVar),
		ClientTwoLeaderXPriv: os.Getenv(ClientTwoLeaderXPrivEnvVar),
	}, nil
}
