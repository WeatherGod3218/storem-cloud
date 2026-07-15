package s3

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client
var bucketName = os.Getenv("AWS_S3_BUCKET")

func CheckExistence(input string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(input),
	})

	logging.Logger.Info(fmt.Sprintf("HeadObject: %v", err))
}

func ListObjects() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := S3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		logging.Logger.Fatal(err)
	}

	for _, obj := range resp.Contents {
		logging.Logger.Infof("S3 object: %s", *obj.Key)
	}
}

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

	S3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String("https://s3.csh.rit.edu/")
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
	})
	logging.Logger.Info("AWS S3 client initialized")
	return nil
}
