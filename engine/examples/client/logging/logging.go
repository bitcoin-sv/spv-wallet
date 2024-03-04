package main

import (
	"context"
	"log"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/logging"
)

func main() {
	client, err := engine.NewClient(
		context.Background(),                          // Set context
		engine.WithLogger(logging.GetDefaultLogger()), // Example of using a custom logger
	)
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	defer func() {
		_ = client.Close(context.Background())
	}()

	log.Println("client loaded!", client.UserAgent())
}
