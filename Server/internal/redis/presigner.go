package redis

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/WeatherGod3218/storem-cloud-server/internal/logging"
	"github.com/WeatherGod3218/storem-cloud-server/internal/s3"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type URLTypes string

const (
	Video     URLTypes = "video"
	Thumbnail URLTypes = "thumbnail"
)

const MAX_TTL = time.Minute * 10
const MIN_TTL = time.Minute * 1

func getRandomTTL() time.Duration {
	ttlRange := MAX_TTL - MIN_TTL
	if ttlRange < 0 {
		logging.Logger.WithFields(logrus.Fields{"max": MAX_TTL, "min": MIN_TTL}).Warn("Invalid Redis Time Configuration!")
		return time.Minute * 2
	}

	randDuration := rand.N(ttlRange)

	return randDuration + MIN_TTL
}

func refreshKeyVal(itemKey string, ver URLTypes) (string, error) {
	var newUrl string
	var err error
	switch ver {
	case Thumbnail:
		newUrl, err = s3.GetThumbnailPresignedURL(itemKey)
	case Video:
		newUrl, err = s3.CreateGetPresignedVideoURL(itemKey)
	default:
		return "", fmt.Errorf("unknown url type %s", ver)
	}
	if err != nil {
		return "", err
	}

	return newUrl, nil
}

func GetPresignedURL(s3Key string, ver URLTypes) (string, error) {
	if !RedisInitalized {
		return "", errors.New("redis was not intialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	cacheKey := fmt.Sprintf("presigned:%s%s", ver, s3Key)

	url, err := client.Get(ctx, cacheKey).Result()
	if err == nil {
		return url, nil
	}
	if err != redis.Nil {
		return "", err
	}

	lockKey := fmt.Sprintf("lock:%s%s", ver, s3Key)

	setLock, err := client.SetNX(ctx, lockKey, 1, time.Second*10).Result()
	if err != nil {
		return "", err
	}

	if setLock {
		defer client.Del(ctx, lockKey)

		newURL, err := refreshKeyVal(s3Key, ver)
		if err != nil {
			return "", err
		}

		ttl := getRandomTTL()

		if err := client.Set(ctx, cacheKey, newURL, ttl).Err(); err != nil {
			return "", err
		}

		return newURL, nil
	}

	//IM BACK ON MY BULLSHITTTTT
	//TODO: Make this better :sob:
	for range 40 {
		time.Sleep(100 * time.Millisecond)

		url, err := client.Get(ctx, cacheKey).Result()
		if err == nil {
			return url, nil
		}
		if err != redis.Nil {
			return "", err
		}
	}

	return "", errors.New("timed out waiting on redis refresh")
}
