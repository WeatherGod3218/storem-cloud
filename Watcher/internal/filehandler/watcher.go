package filehandler

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

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

func ScanSubDirectories(config models.Config, watcher *fsnotify.Watcher, dir string, wg *sync.WaitGroup, recurseAmount int) error {
	if recurseAmount <= 0 {
		return nil
	}

	items, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.IsDir() {
			err := AddDirectory(watcher, dir, item.Name())
			if err != nil {
				logging.Logger.WithFields(logrus.Fields{"error": err, "module": "filehandler", "method": "ScanSubDirectories"}).Warning(fmt.Sprintf("failed to add the directory %s", item.Name()))
				return err
			}
			err = ScanSubDirectories(config, watcher, filepath.Join(dir, item.Name()), wg, recurseAmount-1)
			if err != nil {
				logging.Logger.WithFields(logrus.Fields{"error": err, "module": "filehandler", "method": "ScanSubDirectories"}).Warning(fmt.Sprintf("failed to handle sub directories for %s", item.Name()))
				return err
			}
		} else {
			wg.Add(1)
			go func() {
				defer wg.Done()
				AddFileToVerifyList(filepath.Join(dir, item.Name()))
				if err != nil {
					logging.Logger.WithFields(logrus.Fields{"error": err, "module": "filehandler", "method": "ScanSubDirectories"}).Warning(fmt.Sprintf("failed to hash the item %s", item.Name()))
				}
			}()
		}
	}

	return nil
}

func InitWatcher() (*fsnotify.Watcher, error) {
	confile, err := os.ReadFile("config.toml")
	if err != nil {
		return nil, fmt.Errorf("unable to load config file")
	}

	var config models.Config

	err = toml.Unmarshal(confile, &config)
	if err != nil {
		return nil, fmt.Errorf("error loading config file %w", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error in Creating watcher %w", err)
	}

	if len(config.Directories) == 0 {
		return nil, fmt.Errorf("base Directory was not set")
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

	var wg sync.WaitGroup

	for _, dir := range config.Directories {
		err = watcher.Add(dir)
		if err != nil {
			return nil, err
		}

		err = ScanSubDirectories(config, watcher, dir, &wg, config.Levels)
		if err != nil {
			return nil, err
		}
	}

	wg.Wait()
	return watcher, err
}
