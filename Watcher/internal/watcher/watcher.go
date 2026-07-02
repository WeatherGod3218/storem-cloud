package watcher

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml/v2"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/logging"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/models"
	"github.com/sirupsen/logrus"
)

func AddDirectory(watcher *fsnotify.Watcher, path string, dir string) error {
	err := watcher.Add(filepath.Join(path, dir))
	logging.Logger.WithFields(logrus.Fields{"directory": dir, "module": "watcher", "method": "AddDirectory"}).Info("Added a new directory!")
	if err != nil {
		return err
	}

	return nil
}

func ScanSubDirectories(watcher *fsnotify.Watcher, dir string, recurseAmount int) error {
	if recurseAmount <= 0 {
		return nil
	}

	items, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.IsDir() {
			err = AddDirectory(watcher, dir, item.Name())
			if err != nil {
				return err
			}

			err = ScanSubDirectories(watcher, filepath.Join(dir, item.Name()), recurseAmount-1)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func InitWatcher() (*fsnotify.Watcher, error) {

	confile, err := os.ReadFile("config.toml")
	if err != nil {
		return nil, fmt.Errorf("Unable to Load Config File")
	}

	var config models.Config

	err = toml.Unmarshal(confile, &config)
	if err != nil {
		return nil, fmt.Errorf("Error loading config file!")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("Error in Creating Watcher")
	}

	baseDirectory := config.Directory
	if baseDirectory == "" {
		return nil, fmt.Errorf("Base Directory was not set!")
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Create) {
					logging.Logger.WithFields(logrus.Fields{"module": "watcher", "method": "InitWatcher"}).Info("New File!")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logging.Logger.WithFields(logrus.Fields{"error": err, "module": "watcher", "method": "InitWatcher"}).Warn("Watcher error occured!")
			}
		}
	}()

	err = watcher.Add(baseDirectory)
	if err != nil {
		return nil, err
	}

	err = ScanSubDirectories(watcher, baseDirectory, config.Levels)
	if err != nil {
		return nil, err
	}

	return watcher, err
}
