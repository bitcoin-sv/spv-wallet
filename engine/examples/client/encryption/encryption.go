package main

import (
	"context"
	"log"
	"os"

	"github.com/bitcoin-sv/spv-wallet/engine"
)

func main() {
	client, err := engine.NewClient(
		context.Background(), // Set context
		engine.WithDebugging(),  // Enable debugging (verbose logs)
		engine.WithEncryption(os.Getenv("SPV_WALLET_ENCRYPTION_KEY")), // Encryption key for external public keys (paymail)
	)
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	defer func() {
		_ = client.Close(context.Background())
	}()

	log.Println("client loaded!", client.UserAgent(), "debugging: ", client.IsDebug())
}
