package verify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/WeatherGod3218/storem-cloud-watcher/internal/logging"
	"github.com/WeatherGod3218/storem-cloud-watcher/internal/models"
	"github.com/WeatherGod3218/storem-cloud-watcher/internal/upload"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var filesToVerify map[string]string = make(map[string]string)
var mutex sync.Mutex = sync.Mutex{}

const MINIMUM_VIDEO_LENGTH = 5

func AddFileToVerifyList(config models.Config, file string, baseDir string) {
	mutex.Lock()
	defer mutex.Unlock()

	if !config.IncludeDirectoryPath {
		file = strings.TrimPrefix(file, baseDir)
	}
	filesToVerify[file] = baseDir

}

func GetAllFilesToVerify(config models.Config) []string {
	retList := make([]string, len(filesToVerify))
	index := 0
	for file, baseDir := range filesToVerify {
		if !config.IncludeDirectoryPath {
			retList[index] = strings.TrimPrefix(file, baseDir)
		}
		retList[index] = strings.TrimPrefix(retList[index], string(filepath.Separator))
		index++
	}
	return retList
}

func GetBaseDirFromFile(file string) string {
	file = strings.TrimPrefix(file, string(filepath.Separator)) //incase it has a leading slash
	file = fmt.Sprintf("%s%s", string(filepath.Separator), file)
	return filesToVerify[file]
}

func ValidateFilesForBackup(credentials models.Credentials, config models.Config) error {
	mutex.Lock()
	defer mutex.Unlock()

	fileList := GetAllFilesToVerify(config)

	jsonBytes, err := json.Marshal(fileList)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v2/videos/verify", credentials.ServerURL), bytes.NewBuffer(jsonBytes))
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
		logging.Logger.Infof("Found %d files to upload, against %d total files. %d backed", len(notBacked), len(fileList), (len(fileList) - len(notBacked)))

		for _, file := range notBacked {
			wg.Go(func() error {
				err := upload.UploadVideo(config, file, GetBaseDirFromFile(file))
				if err != nil {
					logging.Logger.WithFields(logrus.Fields{"module": "verify", "method": "ValidateFilesForBackup", "error": err}).Warn("Failure backing up video!")
				}
				return err
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
