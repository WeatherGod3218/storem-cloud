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
	Email        string   `yaml:"email"`
	SubjectID    string   `yaml:"subjectId"`
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

func GetUserByAccessId(id string) (*User, bool) {
	user, ok := UsersById[id]
	return user, ok
}

func GetUserByToken(token any, exists bool) *User {
	if !exists {
		return nil
	}

	userToken, ok := token.(*models.UserToken)
	if !ok || userToken.Email == nil {
		return nil
	}

	userStruct, ok := GetUserByEmail(*userToken.Email)
	if !ok {
		return nil
	}

	if userToken.SubjectId != nil {
		logging.Logger.Infof("Got User! %s %s", *userToken.SubjectId, *userToken.Email)
	}

	userStruct.SubjectID = *userToken.SubjectId
	userStruct.Email = *userToken.Email

	return userStruct
}

func (user *User) CanModifyVideoData() bool {
	switch strings.ToLower(user.Role) {
	case "owner", "admin":
		return true
	default:
		return false
	}
}

// todo ig?
func (user *User) CanViewVideo(video *models.GetVideoDataDatabase) bool {
	switch strings.ToLower(user.Role) {
	case "owner", "admin":
		return true
	default:
		return false
	}
}

func (user *User) CanViewPrivateVideos() bool {
	switch strings.ToLower(user.Role) {
	case "owner", "admin":
		return true
	default:
		return false
	}
}

func (user *User) CanViewAuditLogs() bool {
	switch strings.ToLower(user.Role) {
	case "owner", "admin":
		return true
	default:
		return false
	}
}
