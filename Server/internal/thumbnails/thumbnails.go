package thumbnails

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/WeatherGod3218/weather-reels-server/internal/s3"
)

const THUMBNAIL_SECONDS_BEFORE_END float64 = 5

func RemoveThumbnail(fileName string) error {
	err := os.Remove(fileName)
	return err
}

func GenerateThumbnailFromVideo(uploadId string, fileLength float64) error {
	url, err := s3.CreateGetPresignedVideoURL(uploadId)
	if err != nil {
		return fmt.Errorf("unable to get presigned video url %s", err)
	}

	var screenshotTime string
	if fileLength > THUMBNAIL_SECONDS_BEFORE_END*1.5 { //Bigger so you dont get a 0.01 thumbnail or somethin
		screenshotTime = fmt.Sprintf("%f", fileLength-THUMBNAIL_SECONDS_BEFORE_END)
	} else {
		screenshotTime = fmt.Sprintf("%f", fileLength/2)
	}
	cmd := exec.Command(
		"ffmpeg",
		"-ss", screenshotTime,
		"-i", url,
		"-frames:v", "1",
		"-f", "image2",
		"-vcodec", "png",
		"pipe:1",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to run ffmpeg command %s, %s", err, stderr.String())
	}

	imageBytes := stdout.Bytes()
	if len(imageBytes) == 0 {
		return fmt.Errorf("no output from thumbnail generation %s", stderr.String())
	}

	if err := s3.StoreThumbnailImageBytes(uploadId, imageBytes); err != nil {
		return err
	}
	return nil
}
