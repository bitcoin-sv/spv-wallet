package main

import (
	"context"
	"log"

	"github.com/bitcoin-sv/spv-wallet/engine"
)

func main() {
	client, err := engine.NewClient(
		context.Background(), // Set context
		engine.WithDebugging(),  // Enable debugging (verbose logs)
		engine.WithChainstateOptions(true, true, true, true), // Broadcasting enabled by default
	)
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	defer func() {
		_ = client.Close(context.Background())
	}()

	log.Println("client loaded!", client.UserAgent(), "debugging: ", client.IsDebug())
}
