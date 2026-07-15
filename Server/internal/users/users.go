package users

import (
	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/models"
)

var Users map[string]models.UserConfig = make(map[string]models.UserConfig)

func InitUsers(config models.Config) {
	for _, user := range config.Users {
		Users[user.UserId] = user
		logging.Logger.Infof("Initalized User %s!", user.DisplayName)
	}
}

func VerifyUser(userId string) bool {
	_, ok := Users[userId]
	return ok
}

func GetUsername(userId string) string {
	return Users[userId].DisplayName
}
