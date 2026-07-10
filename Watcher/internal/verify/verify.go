package verify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/logging"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/upload"
)

var filesToVerify []string = make([]string, 0)
var mutex sync.Mutex = sync.Mutex{}

const MINIMUM_VIDEO_LENGTH = 5

func AddFileToVerifyList(file string) {
	mutex.Lock()
	defer mutex.Unlock()

	filesToVerify = append(filesToVerify, file)
}

func ValidateFilesForBackup() error {
	mutex.Lock()
	defer mutex.Unlock()

	verifyList := filesToVerify

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
		err := upload.UploadVideo(notBacked[0])
		if err != nil {
			logging.Logger.Info(fmt.Sprintf("Failure in backup: %s", err))
		}
	}
	return nil
}
