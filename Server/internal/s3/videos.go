package s3

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const PRESIGNED_URL_TIME = 15

func CreatePermanentUrl(video string) (string, error) {
	bucket := os.Getenv("AWS_S3_UID")
	if bucket == "" {
		return "", fmt.Errorf("AWS_S3_UID environment variable is not set")
	}
	region := GetRegion()
	key := CreateVideoKey(video)

	fileUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
	return fileUrl, nil
}

func CreatePresignedUrl(video string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	svc := s3.NewFromConfig(cfg)
	presign := s3.NewPresignClient(svc)

	bucket := os.Getenv("AWS_S3_UID")
	if bucket == "" {
		return "", fmt.Errorf("AWS_S3_UID environment variable is not set")
	}
	key := CreateVideoKey(video)

	lifetime := 15 * time.Minute

	request, err := presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(lifetime))
	if err != nil {
		return "", fmt.Errorf("unable to create presigned url")
	}

	return request.URL, nil
}
