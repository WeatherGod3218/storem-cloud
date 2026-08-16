package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/WeatherGod3218/storem-cloud-server/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const PRESIGNED_URL_TIME = 15

const (
	MIN_PART_SIZE = 5 * 1024 * 1024
	MAX_PART_SIZE = 5 * 1024 * 1024 * 1024
	MAX_PARTS     = 10_000
	TARGET_PARTS  = 32
)

func calculateParts(fileSize int64) ([]models.VideoPart, error) {
	if fileSize <= 0 {
		return nil, fmt.Errorf("invalid file size: %d", fileSize)
	}

	if fileSize < MIN_PART_SIZE {
		return []models.VideoPart{{PartNumber: 1, Offset: 0, Size: fileSize}}, nil
	}

	partSize := ceilDiv(fileSize, TARGET_PARTS)

	if partSize < MIN_PART_SIZE {
		partSize = MIN_PART_SIZE
	}

	numParts := ceilDiv(fileSize, partSize)

	if numParts > MAX_PARTS {
		partSize = ceilDiv(fileSize, MAX_PARTS)
		numParts = ceilDiv(fileSize, partSize)
	}

	if partSize > MAX_PART_SIZE {
		return nil, fmt.Errorf("file too large for multipart upload: %d bytes exceeds max object size", fileSize)
	}

	var offset int64

	parts := make([]models.VideoPart, 0, numParts)
	partNumber := int32(1)

	for offset < fileSize {
		size := partSize
		remaining := fileSize - offset
		if remaining < size {
			size = remaining
		}
		parts = append(parts, models.VideoPart{
			PartNumber: partNumber,
			Offset:     offset,
			Size:       size,
		})
		offset += size
		partNumber++
	}

	return parts, nil
}

func ceilDiv(a, b int64) int64 {
	return (a + b - 1) / b
}

func CreatePermanentVideoURL(uploadId string) (string, error) {
	if BucketName == "" {
		return "", fmt.Errorf("bucket environment variable is not set")
	}
	region := GetRegion()
	key := uploadId

	fileUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", BucketName, region, key)
	return fileUrl, nil
}

func DeleteVideo(uploadId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if BucketName == "" {
		return fmt.Errorf("bucket environment variable is not set")
	}

	key := uploadId

	_, err := S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to create presigned url")
	}

	return nil
}

func CreateGetPresignedVideoURL(uploadId string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	presign := s3.NewPresignClient(S3Client)

	if BucketName == "" {
		return "", fmt.Errorf("bucket environment variable is not set")
	}

	key := uploadId

	lifetime := GetPresignURLTime()

	request, err := presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(lifetime))
	if err != nil {
		return "", fmt.Errorf("unable to create presigned url")
	}

	return request.URL, nil
}
