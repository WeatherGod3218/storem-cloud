package thumbnails

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RemoveThumbnail(fileName string) error {
	err := os.Remove(fileName)
	return err
}

func GenerateThumbnailFromVideo(fileName string) (string, error) {
	base := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	outPath := fmt.Sprintf("thumbnails/%s.png", base)

	fileCmd := exec.Command("ffmpeg", "-ss", "00:00:05", "-i", fmt.Sprintf("%s", fileName), "-vframes", "1", outPath)

	_, err := fileCmd.Output()
	if err != nil {
		return "", fmt.Errorf("unable to run ffmpeg command")
	}
	return outPath, nil
}
