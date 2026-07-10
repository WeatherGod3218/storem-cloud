package backup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/filehandler"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/logging"
	"github.com/alfg/mp4"
)

const MINIMUM_VIDEO_LENGTH = 5

func getMP4Length(fileName string) (float64, error) {
	file, err := mp4.Open(fileName)
	if err != nil {
		return 0, err
	}

	if file.Moov == nil || file.Moov.Mvhd == nil {
		return 0, fmt.Errorf("file does not contain metadata")
	}

	durationSeconds := (float64(file.Moov.Mvhd.Duration) / float64(file.Moov.Mvhd.Timescale))
	return durationSeconds, nil
}

func ValidateFilesForBackup() error {
	verifyList := filehandler.GetVerifyList()

	logging.Logger.Info("Hashlist?")
	jsonBytes, err := json.Marshal(verifyList)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v1/videos/verify", os.Getenv("SERVER_URL")), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to start video request!")
	}

	var notBacked []string
	err = json.NewDecoder(res.Body).Decode(&notBacked)
	if err != nil {
		return err
	}

	if len(notBacked) > 0 {
		err := BackupFile(notBacked[0])
		if err != nil {
			logging.Logger.Info(fmt.Sprintf("Failure in backup: %s", err))
		}
	}
	return nil
}
