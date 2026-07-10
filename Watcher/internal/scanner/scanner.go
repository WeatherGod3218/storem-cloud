package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/models"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/verify"
)

func ScanDirectory(config models.Config, wg *sync.WaitGroup, path string, recurseAmount int, includeSubDirs bool, watcher *fsnotify.Watcher) error {
	if recurseAmount == 0 {
		return nil
	}

	items, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.IsDir() {
			newPath := filepath.Join(path, item.Name())

			if config.BackupDuringRuntime {
				if err := watcher.Add(newPath); err != nil {
					return fmt.Errorf("failed to add watcher to directory %s with error %s", newPath, err)
				}
			}

			if includeSubDirs {
				if err := ScanDirectory(config, wg, newPath, recurseAmount-1, false, watcher); err != nil {
					return fmt.Errorf("failed to scan the directory %s with error %s", newPath, err)
				}
			}
		} else {
			wg.Add(1)
			go func() {
				defer wg.Done()
				verify.AddFileToVerifyList(filepath.Join(path, item.Name()))
			}()
		}
	}

	return nil
}

func ScanFiles(config models.Config, watcher *fsnotify.Watcher) error {

	var wg sync.WaitGroup

	for _, dir := range config.Directories {
		if err := ScanDirectory(config, &wg, dir.Path, dir.SubDirectoryLevels, dir.IncludeSubDirectories, watcher); err != nil {
			return err
		}
	}

	wg.Wait()
	return nil
}
