package backup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/filehandler"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/logging"
)

func ValidateHashedFiles() error {
	hashList := filehandler.GetHashedFilesList()

	logging.Logger.Info("Hashlist?")
	jsonBytes, err := json.Marshal(hashList)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/videos/verify", os.Getenv("SERVER_URL")), bytes.NewBuffer(jsonBytes))
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

	var notBacked []string
	err = json.NewDecoder(res.Body).Decode(&notBacked)
	if err != nil {
		return err
	}

	for _, file := range notBacked {
		logging.Logger.Info(fmt.Sprintf("Needs to be validated: %s", file))
	}
	return nil
}
