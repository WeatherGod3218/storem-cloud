package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/logging"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/watcher"
	"github.com/sirupsen/logrus"
)

func main() {
	watcher, err := watcher.InitWatcher()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Failed to start watcher!")
	}
	defer watcher.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM) // my probe returns

	logging.Logger.Info("Watcher running, waiting for events...")
	<-sigChan

	logging.Logger.Info("Shutting down watcher...")
}
