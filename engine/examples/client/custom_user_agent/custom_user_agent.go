package main

import (
	"context"
	"log"

	"github.com/bitcoin-sv/spv-wallet/engine"
)

func main() {
	client, err := engine.NewClient(
		context.Background(),                      // Set context
		engine.WithUserAgent("my-custom-user-agent"), // Custom user agent
	)
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	defer func() {
		_ = client.Close(context.Background())
	}()

	log.Println("client loaded!", client.UserAgent())
}
