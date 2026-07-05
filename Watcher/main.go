package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/backup"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/filehandler"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/logging"
	"github.com/sirupsen/logrus"
)

func main() {
	watcher, err := filehandler.InitWatcher() // Make Watcher
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Failed to start watcher!")
	}
	defer watcher.Close()

	//Hash Files

	//Send Hashes Files to Server For Validation

	time.Sleep(5 * time.Second)

	err = backup.ValidateHashedFiles()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Warning("Failed to send data!")
	}
	// Backup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM) // my probe returns

	logging.Logger.Info("Watcher running, waiting for events...")
	<-sigChan

	logging.Logger.Info("Shutting down watcher...")
}
