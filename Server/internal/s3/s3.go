package s3

import (
	"context"
	"time"

	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var cfg aws.Config

func InitS3() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(GetRegion()),
	)
	if err != nil {
		return err
	}

	logging.Logger.Info("S3 Connection has been connected!")
	cfg = awsCfg
	return nil
}
