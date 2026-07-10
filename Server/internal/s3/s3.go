package s3

import (
	"context"
	"time"

	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client

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

	s3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String("https://s3.csh.rit.edu/")
	})
	logging.Logger.Info("AWS S3 client initialized")
	return nil
}
