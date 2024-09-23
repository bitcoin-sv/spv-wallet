package main

import (
	"context"
	"log"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine"
)

func main() {
	client, err := engine.NewClient(
		context.Background(), // Set context
		engine.WithCronCustomPeriod(engine.CronJobNameDraftTransactionCleanUp, 2*time.Second),
		engine.WithCronCustomPeriod(engine.CronJobNameSyncTransaction, 4*time.Second),
	)
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	defer func() {
		_ = client.Close(context.Background())
	}()

	// wait for the customized cron jobs to run at least once
	time.Sleep(8 * time.Second)

	log.Println("client loaded!", client.UserAgent())
}
