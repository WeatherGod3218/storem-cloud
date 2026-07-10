package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/backup"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/filehandler"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/logging"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/models"
	"github.com/go-yaml/yaml"
	"github.com/sirupsen/logrus"
)

func main() {

	watcher, err := filehandler.InitWatcher() // Make Watcher
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Failed to start watcher!")
	}
	defer watcher.Close()

	cfgFile, err := os.ReadFile("config.yaml")
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Failed to start watcher!")
	}

	var config models.Config
	err = yaml.Unmarshal(cfgFile, &config)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Error loading config file!")
	}
	//Send Hashes Files to Server For Validation

	time.Sleep(5 * time.Second)

	err = backup.ValidateFilesForBackup()
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
