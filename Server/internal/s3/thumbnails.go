package s3

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func createThumbnailKey(video string) string {
	return url.PathEscape(fmt.Sprintf("keys/%s", video))
}

func CreatePermanentThumbnailURL(video string) (string, error) {
	bucket := os.Getenv("AWS_S3_UID")
	if bucket == "" {
		return "", fmt.Errorf("AWS_S3_UID environment variable is not set")
	}
	region := GetRegion()
	key := createThumbnailKey(video)

	fileUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
	return fileUrl, nil
}

func CreatePresignedThumbnailURL(video string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	presign := s3.NewPresignClient(s3Client)

	bucket := os.Getenv("AWS_S3_UID")
	if bucket == "" {
		return "", fmt.Errorf("AWS_S3_UID environment variable is not set")
	}
	key := createThumbnailKey(video)

	lifetime := GetPresignURLTime()

	request, err := presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(lifetime))
	if err != nil {
		return "", fmt.Errorf("unable to create presigned url")
	}

	return request.URL, nil
}
