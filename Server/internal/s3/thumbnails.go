package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func convertToThumbnailKey(fileName string) string {
	return fmt.Sprintf("thumbnails/%s.png", fileName)
}

func CreatePermanentThumbnailURL(video string) (string, error) {
	bucket := bucketName
	if bucket == "" {
		return "", fmt.Errorf("AWS_S3_UID environment variable is not set")
	}
	region := GetRegion()
	key := convertToThumbnailKey(video)

	fileUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
	return fileUrl, nil
}

func StoreThumbnailImageBytes(uploadId string, byteData []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bucket := bucketName
	if bucket == "" {
		return errors.New("AWS_S3_UID environment variable is not set")
	}
	key := convertToThumbnailKey(uploadId)

	_, err := S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(byteData),
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		return fmt.Errorf("unable to upload thumbnail: %w", err)
	}

	return nil
}

func GetThumbnailPresignedURL(id string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	presign := s3.NewPresignClient(S3Client)

	bucket := bucketName
	if bucket == "" {
		return "", fmt.Errorf("AWS_S3_UID environment variable is not set")
	}
	key := convertToThumbnailKey(id)

	lifetime := GetPresignURLTime()

	request, err := presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(lifetime))
	if err != nil {
		return "", fmt.Errorf("unable to create presigned url")
	}

	return request.URL, nil
}
