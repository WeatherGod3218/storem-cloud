package s3

import (
	"fmt"
	"net/url"
	"os"
)

func GetRegion() string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return "us-east-1"
	}
	return region
}

func CreateVideoKey(video string) string {
	return url.PathEscape(fmt.Sprintf("videos/%s", video))
}
