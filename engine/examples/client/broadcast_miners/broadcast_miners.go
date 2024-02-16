package main

import (
	"context"
	"log"
	"os"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	broadcastclient "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
	"github.com/bitcoin-sv/spv-wallet/engine/logging"
)

func buildBroadcastClient() broadcast.Client {
	logger := logging.GetDefaultLogger()
	builder := broadcastclient.Builder().WithArc(
		broadcastclient.ArcClientConfig{
			APIUrl: "https://tapi.taal.com/arc",
			Token:  os.Getenv("SPV_WALLET_TAAL_API_KEY"),
		},
		logger,
	)

	return builder.Build()
}

func main() {
	ctx := context.Background()

	client, err := engine.NewClient(
		ctx,
		engine.WithBroadcastClient(buildBroadcastClient()),
	)
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	defer client.Close(ctx)

	log.Println("client loaded!", client.UserAgent())
}
