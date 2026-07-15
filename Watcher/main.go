package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/logging"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/models"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/scanner"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/upload"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/verify"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/watcher"

	"github.com/fsnotify/fsnotify"
	"github.com/go-yaml/yaml"
	"github.com/sirupsen/logrus"
)

func main() {
	cfgFile, err := os.ReadFile("config.yaml")
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Failed read in config file!")
	}

	var config models.Config
	err = yaml.Unmarshal(cfgFile, &config)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Error loading config file!")
	}

	credFile, err := os.ReadFile("credentials.yaml")
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Failed to read in cred file!")
	}

	var credentials models.Credentials
	err = yaml.Unmarshal(credFile, &credentials)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Error loading credentials file!")
	}

	if len(config.Directories) <= 0 {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("No directories were set to be backed up!")
	}

	upload.InitTusio(credentials)

	var wtch *fsnotify.Watcher

	if config.BackupDuringRuntime {
		wtch, err = watcher.InitWatcher()
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Failed to start watcher!")
		}
		defer wtch.Close()
	}

	if err := scanner.ScanFiles(config, wtch); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Error scanning files for verification")
	}

	time.Sleep(time.Second * 5)

	if err := verify.ValidateFilesForBackup(credentials, config); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Error validating files for backup!")
	}

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Warning("Failed to send data!")
	}
	// Backup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM) // my probe returns

	logging.Logger.Info("Watcher running, waiting for events...")

	// //TESTING PURPOSES

	// firstPage := models.GetVideoGroupRequest{
	// 	Timestamp: nil,
	// 	RowId:     "",
	// }
	// pageBytes, _ := json.Marshal(firstPage)

	// req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/videos/group", os.Getenv("SERVER_URL")), bytes.NewBuffer(pageBytes))
	// if err != nil {
	// 	logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Failed to create main req")
	// }
	// req.Header.Set("Content-Type", "application/json")

	// client := &http.Client{}
	// completeRes, err := client.Do(req)
	// if err != nil {
	// 	logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Failed to do request")
	// }
	// defer completeRes.Body.Close()
	// if completeRes.StatusCode != http.StatusOK {
	// 	logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Failed to generate video id")
	// }

	// var response models.GetVideoGroupResponse
	// bodyBytes, err := io.ReadAll(completeRes.Body)
	// if err != nil {
	// 	logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Failed to read")
	// }

	// logging.Logger.Infof(
	// 	"Status: %d, Body: %q",
	// 	completeRes.StatusCode,
	// 	string(bodyBytes),
	// )

	// err = json.Unmarshal(bodyBytes, &response)
	// if err != nil {
	// 	logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Failed to unmarshal")
	// }

	// for _, part := range response.Videos {
	// 	logging.Logger.Info(part.VideoURL)
	// }

	<-sigChan

	logging.Logger.Info("Shutting down watcher...")
}
