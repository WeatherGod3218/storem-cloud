package verify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"math/rand/v2"
	"net/http"
	"os"
	"slices"
	"sync"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/logging"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/models"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/upload"
)

var filesToVerify map[string]string = make(map[string]string)
var mutex sync.Mutex = sync.Mutex{}

const MINIMUM_VIDEO_LENGTH = 5

func AddFileToVerifyList(file string, baseDir string) {
	mutex.Lock()
	defer mutex.Unlock()

	filesToVerify[file] = baseDir
}

func ValidateFilesForBackup(credentials models.Credentials, config models.Config) error {
	mutex.Lock()
	defer mutex.Unlock()

	jsonBytes, err := json.Marshal(slices.Collect(maps.Keys(filesToVerify)))
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v2/videos/verify", os.Getenv("SERVER_URL")), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials.ServerAccessCode))

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
		randNum := rand.IntN(len(notBacked))
		fileToBack := notBacked[randNum]
		err := upload.UploadVideo(config, fileToBack, filesToVerify[fileToBack])
		if err != nil {
			logging.Logger.Info(fmt.Sprintf("Failure in backup: %s", err))
		}
	}
	return nil
}
