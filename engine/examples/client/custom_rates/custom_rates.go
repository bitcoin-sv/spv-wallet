package main

import (
	"context"
	"log"
	"os"
	"time"

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
	const testXPub = "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J"

	client, err := engine.NewClient(
		ctx,
		engine.WithAutoMigrate(engine.BaseModels...),
		engine.WithBroadcastClient(buildBroadcastClient()),
	)
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	defer client.Close(ctx)

	xpub, err := client.NewXpub(ctx, testXPub)
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	draft, err := client.NewTransaction(ctx, xpub.RawXpub(), &engine.TransactionConfig{
		ExpiresIn: 10 * time.Second,
		SendAllTo: &engine.TransactionOutput{To: os.Getenv("SPV_WALLET_MY_PAYMAIL")},
	})
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	// Custom fee
	log.Println("fee unit", draft.Configuration.FeeUnit)
}
