package main

import (
	"context"
	"log"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	// EXAMPLE: new relic application
	// replace this with your ALREADY EXISTING new relic application
	app, err := tester.GetNewRelicApp("test-app")
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	var client engine.ClientInterface
	client, err = engine.NewClient(
		newrelic.NewContext(context.Background(), app.StartTransaction("test-txn")), // Set context
		engine.WithNewRelic(app), // New relic application (from your own application or server)
	)
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	defer func() {
		_ = client.Close(context.Background())
	}()

	log.Println("client loaded!", client.UserAgent())
}
