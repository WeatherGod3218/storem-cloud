package redis

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client
var RedisInitalized bool = false

func InitRedis() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return err
	}

	newClient := redis.NewClient(opt)
	_, err = newClient.Ping(ctx).Result()
	if err != nil {
		_ = newClient.Close()
		return err
	}

	RedisInitalized = true
	client = newClient
	return nil
}
