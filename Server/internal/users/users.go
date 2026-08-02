package users

import (
	"strings"

	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/models"
)

type User struct {
	UserID       string   `yaml:"userId"`
	Emails       []string `yaml:"emails"`
	DisplayName  string   `yaml:"displayName"`
	TotalStorage int      `yaml:"totalStorageGB"`
	Role         string   `yaml:"role"`
}

var UsersById map[string]*User = make(map[string]*User)
var UsersByEmail map[string]*User = make(map[string]*User)

func InitUsers(config models.Config) {
	for _, userConfig := range config.Users {
		user := &User{
			UserID:       userConfig.UserID,
			Emails:       userConfig.Emails,
			DisplayName:  userConfig.DisplayName,
			TotalStorage: userConfig.TotalStorage,
			Role:         userConfig.Role,
		}

		UsersById[user.UserID] = user
		for _, email := range user.Emails {
			UsersByEmail[strings.ToLower(email)] = user
		}
		logging.Logger.Infof("Initalized User %s!", user.DisplayName)
	}
}

func VerifyUser(userId string) bool {
	_, ok := UsersById[userId]
	return ok
}

func GetUsername(userId string) string {
	user := *UsersById[userId]
	return user.DisplayName
}

func GetUserByEmail(email string) (*User, bool) {
	user, ok := UsersByEmail[strings.ToLower(email)]
	return user, ok
}

func GetUserByToken(token any, exists bool) *User {
	if !exists {
		logging.Logger.Info("Doesnt exist")
		return nil
	}

	userToken, ok := token.(*models.UserToken)
	if !ok || userToken.Email == nil {
		logging.Logger.Info("Not okay")
		return nil
	}

	userStruct, ok := GetUserByEmail(*userToken.Email)
	if !ok {
		if userToken.Email != nil {
			logging.Logger.Infof("didnt verify email: %s", *userToken.Email)
		}

		logging.Logger.Info("was no email")

		return nil
	}

	if userToken.SubjectId != nil {
		logging.Logger.Infof("Got User! %s %s", *userToken.SubjectId, *userToken.Email)
	}

	return userStruct
}

func (user *User) CanModifyVideoData() bool {
	switch user.Role {
	case "owner", "admin":
		return true
	default:
		return false
	}
}
