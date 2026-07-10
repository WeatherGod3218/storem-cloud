package watcher

import (
	"fmt"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/logging"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

func InitWatcher() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error in creating watcher %w", err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Create) {
					logging.Logger.WithFields(logrus.Fields{"module": "watcher", "method": "InitWatcher"}).Info("new file")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logging.Logger.WithFields(logrus.Fields{"error": err, "module": "watcher", "method": "InitWatcher"}).Warn("watcher error occured")
			}
		}
	}()

	return watcher, nil
}
