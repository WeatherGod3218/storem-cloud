package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/s3"
	"github.com/redis/go-redis/v9"
)

type URLTypes string

const (
	Video     URLTypes = "video"
	Thumbnail URLTypes = "thumbnail"
)

func GetPresignedURL(s3Key string, ver URLTypes) (string, error) {
	if !RedisInitalized {
		return "", errors.New("redis was not intialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	cacheKey := fmt.Sprintf("presigned:%s%s", ver, s3Key)
	logging.Logger.Info(cacheKey)

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

		var newUrl string
		var err error
		switch ver {
		case Thumbnail:
			newUrl, err = s3.GetThumbnailPresignedURL(s3Key)
		case Video:
			newUrl, err = s3.CreateGetPresignedVideoURL(s3Key)
		default:
			return "", fmt.Errorf("unknown url type %s", ver)
		}
		if err != nil {
			return "", err
		}

		ttl := s3.GetPresignURLTime() - 5*time.Minute
		if ttl <= 0 {
			ttl = s3.GetPresignURLTime()
		}

		if err := client.Set(ctx, cacheKey, newUrl, ttl).Err(); err != nil {
			return "", err
		}

		return newUrl, nil
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
