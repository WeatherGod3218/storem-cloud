package verify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"slices"
	"sync"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/logging"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/models"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/upload"
	"golang.org/x/sync/errgroup"
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

	var wg errgroup.Group
	wg.SetLimit(upload.MAX_VIDEO_UPLOADS)

	if len(notBacked) > 0 {
		logging.Logger.Infof("Found %d files to upload", len(notBacked))

		for _, file := range notBacked {
			wg.Go(func() error {
				return upload.UploadVideo(config, file, filesToVerify[file])
			})
		}

		if err := wg.Wait(); err != nil {
			return err
		}
	} else {
		logging.Logger.Info("No files to upload!")
	}
	return nil
}
