package s3

import (
	"os"
	"time"
)

func GetRegion() string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return "us-east-1"
	}
	return region
}

func GetPresignURLTime() time.Duration {
	return (15 * time.Minute)
}
