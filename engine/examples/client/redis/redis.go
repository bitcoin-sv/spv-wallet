package main

import (
	"context"
	"log"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/mrz1836/go-cachestore"
)

func main() {
	redisURL := "localhost:6379"
	client, err := engine.NewClient(
		context.Background(), // Set context
		engine.WithRedis(&cachestore.RedisConfig{URL: redisURL}), // Cache
		engine.WithTaskqConfig( // Tasks
			taskmanager.DefaultTaskQConfig("example_queue", taskmanager.WithRedis(redisURL)),
		),
	)
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	defer func() {
		_ = client.Close(context.Background())
	}()

	log.Println("client loaded!", client.UserAgent())
}
