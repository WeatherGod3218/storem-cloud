package models

type DirectoryConfig struct {
	Path                  string `yaml:"path"`
	IncludeSubDirectories bool   `yaml:"includeSubDirectories"`
	SubDirectoryLevels    int    `yaml:"subDirectoryLevels"`
}

type Config struct {
	Directories          []DirectoryConfig `yaml:"directories"`
	HashFullFile         bool              `yaml:"hashFullFile"`
	IncludeDirectoryPath bool              `yaml:"includeDirectoryPath"`
	BackupDuringRuntime  bool              `yaml:"backupDuringRuntime"`
}
