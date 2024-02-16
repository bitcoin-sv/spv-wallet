package tester

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mrz1836/go-cache"
	"github.com/rafaeljusto/redigomock"
)

// LoadMockRedis will load a mocked redis connection
func LoadMockRedis(
	idleTimeout time.Duration,
	maxConnTime time.Duration,
	maxActive int,
	maxIdle int,
) (client *cache.Client, conn *redigomock.Conn) {
	conn = redigomock.NewConn()
	client = &cache.Client{
		DependencyScriptSha: "",
		Pool: &redis.Pool{
			Dial:            func() (redis.Conn, error) { return conn, nil },
			IdleTimeout:     idleTimeout,
			MaxActive:       maxActive,
			MaxConnLifetime: maxConnTime,
			MaxIdle:         maxIdle,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, doErr := c.Do(cache.PingCommand)
				return doErr
			},
		},
		ScriptsLoaded: nil,
	}
	return
}

// LoadRealRedis will load a real redis connection
func LoadRealRedis(
	connectionURL string,
	idleTimeout time.Duration,
	maxConnTime time.Duration,
	maxActive int,
	maxIdle int,
	dependency bool,
	newRelic bool,
) (client *cache.Client, conn redis.Conn, err error) {
	client, err = cache.Connect(
		context.Background(),
		connectionURL,
		maxActive,
		maxIdle,
		maxConnTime,
		idleTimeout,
		dependency,
		newRelic,
	)
	if err != nil {
		return
	}

	conn, err = client.GetConnectionWithContext(context.Background())
	return
}
