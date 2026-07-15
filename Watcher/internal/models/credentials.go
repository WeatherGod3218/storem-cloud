package models

type Credentials struct {
	ServerURL        string `yaml:"serverUrl"`
	ServerAccessCode string `yaml:"serverAccessCode"`
	UserId           string `yaml:"userId"`
}
