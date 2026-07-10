package models

type DirectoryConfig struct {
	Path                  string `yaml:"path"`
	IncludeSubDirectories bool   `yaml:"includeSubdirectories"`
	SubDirectoryLevels    int    `yaml:"subdirectoryLevels"`
}

type Config struct {
	Directories          []DirectoryConfig `yaml:"directories"`
	HashFullFile         bool              `yaml:"hashFullFile"`
	IncludeDirectoryPath bool              `yaml:"includeDirectoryPath"`
	BackupDuringRuntime  bool              `yaml:"backupDuringRuntime"`
}
