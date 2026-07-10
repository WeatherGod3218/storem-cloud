package s3

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const PRESIGNED_URL_TIME = 15

func createVideoKey(video string) string {
	return url.PathEscape(fmt.Sprintf("videos/%s", video))
}

const (
	MIN_PART_SIZE = 5 * 1024 * 1024
	MAX_PART_SIZE = 5 * 1024 * 1024 * 1024
	MAX_PARTS     = 10_000
	TARGET_PARTS  = 32
)

var bucketName = os.Getenv("AWS_S3_BUCKET")

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

func CreatePermanentVideoURL(fileName string) (string, error) {
	bucket := bucketName
	if bucket == "" {
		return "", fmt.Errorf("bucket environment variable is not set")
	}
	region := GetRegion()
	key := createVideoKey(fileName)

	fileUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
	return fileUrl, nil
}

func CreateGetPresignedVideoURL(fileName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	presign := s3.NewPresignClient(s3Client)

	bucket := bucketName
	if bucket == "" {
		return "", fmt.Errorf("bucket environment variable is not set")
	}

	key := createVideoKey(fileName)

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

func CreatePutPresignedVideoURL(fileName string, filesize int64) ([]models.VideoURLPart, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	presign := s3.NewPresignClient(s3Client)

	bucket := bucketName
	if bucket == "" {
		return nil, "", fmt.Errorf("bucket environment variable is not set")
	}
	key := createVideoKey(fileName)

	lifetime := GetPresignURLTime()

	createResp, err := s3Client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, "", err
	}

	parts, err := calculateParts(filesize)
	if err != nil {
		return nil, "", err
	}

	uploadId := *createResp.UploadId

	requests := make([]models.VideoURLPart, len(parts))

	for i, part := range parts {
		req, err := presign.PresignUploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(bucket),
			Key:        aws.String(key),
			UploadId:   aws.String(uploadId),
			PartNumber: aws.Int32(int32(part.PartNumber)),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = lifetime
		})
		if err != nil {
			return nil, "", err
		}

		requests[i] = models.VideoURLPart{
			RequestURL: req.URL,
			Size:       part.Size,
			PartNumber: part.PartNumber,
			Offset:     part.Offset,
		}
	}

	return requests, uploadId, nil
}

func AbortMultipartUpload(fileName string, uploadId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	bucket := bucketName
	if bucket == "" {
		return fmt.Errorf("bucket environment variable is not set")
	}

	key := createVideoKey(fileName)

	input := &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadId),
	}

	_, err := s3Client.AbortMultipartUpload(ctx, input)

	return err
}

func CompleteMultipartUpload(fileName string, uploadId string, completedParts []models.VideoCompletedPart) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	bucket := bucketName
	if bucket == "" {
		return fmt.Errorf("bucket environment variable is not set")
	}

	key := createVideoKey(fileName)

	parts := make([]types.CompletedPart, 0, len(completedParts))

	for _, p := range completedParts {
		parts = append(parts, types.CompletedPart{
			PartNumber: aws.Int32(int32(p.PartNumber)),
			ETag:       aws.String(p.ETag),
		})
	}

	input := &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadId),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: parts,
		},
	}

	_, err := s3Client.CompleteMultipartUpload(ctx, input)

	return err
}
