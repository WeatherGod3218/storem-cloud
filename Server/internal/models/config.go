package models

type UserConfig struct {
	UserID       string `yaml:"userId"`
	DisplayName  string `yaml:"displayName"`
	TotalStorage int    `yaml:"totalStorageGB"`
}
type Config struct {
	Users []UserConfig `yaml:"users"`
}
