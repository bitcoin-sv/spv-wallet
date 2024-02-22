package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/mrz1836/go-datastore"
)

func main() {
	defaultTimeouts := 10 * time.Second

	client, err := engine.NewClient(
		context.Background(), // Set context
		engine.WithSQL(datastore.MySQL, &datastore.SQLConfig{ // Load using a MySQL configuration
			CommonConfig: datastore.CommonConfig{
				Debug:                 true,
				MaxConnectionIdleTime: defaultTimeouts,
				MaxConnectionTime:     defaultTimeouts,
				MaxIdleConnections:    10,
				MaxOpenConnections:    10,
				TablePrefix:           "spv",
			},
			Driver:    datastore.MySQL.String(),
			Host:      "localhost",
			Name:      os.Getenv("DB_NAME"),
			Password:  os.Getenv("DB_PASSWORD"),
			Port:      "3306",
			TimeZone:  "UTC",
			TxTimeout: defaultTimeouts,
			User:      os.Getenv("DB_USER"),
		}),
		engine.WithPaymailSupport([]string{"test.com"}, "example@test.com", false, false),
		engine.WithAutoMigrate(engine.BaseModels...),
	)
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	defer func() {
		_ = client.Close(context.Background())
	}()

	log.Println("client loaded!", client.UserAgent())
}
