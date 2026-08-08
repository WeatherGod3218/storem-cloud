package database

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/WeatherGod3218/weather-reels-server/internal/users"
	"github.com/sirupsen/logrus"
)

type LoggingInserts string

const (
	TagName   LoggingInserts = "<--TAG:(ID)-->"
	VideoName LoggingInserts = "<--VIDEO:(ID)-->"
)

const MAX_ACTION_AMOUNT int = 3
const MAX_ACTIVE_LOG_WRITES int = 12

var sem = make(chan struct{}, MAX_ACTIVE_LOG_WRITES)

var loggingInsertRe = regexp.MustCompile(`<--(TAG|VIDEO):\(([^)]+)\)-->`)

func resolveInserts(input string) string {
	return loggingInsertRe.ReplaceAllStringFunc(input, func(match string) string {
		groups := loggingInsertRe.FindStringSubmatch(match)
		itemType, id := groups[1], groups[2]

		switch itemType {
		case "TAG":
			name, err := GetTagName(id)
			if err != nil {
				name = ""
			}
			return fmt.Sprintf(`"%s" (%s)`, name, id)
		case "VIDEO":
			name, err := GetVideoName(id)
			if err != nil {
				name = ""
			}
			return fmt.Sprintf(`"%s" (%s)`, name, id)
		default:
			return match
		}
	})
}

func ActionTagName(ID string) string {
	return fmt.Sprintf("<--TAG:(%s)-->", ID)
}

func ActionVideoName(ID string) string {
	return fmt.Sprintf("<--VIDEO:(%s)-->", ID)
}

func LogAction(user *users.User, action string) {
	sem <- struct{}{}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		defer func() { <-sem }()

		rowId, err := GenerateUUID()
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "user": *user, "action": action}).Warn("Failed to log event!")
			return
		}

		baseQuery := `
		INSERT INTO actions (row_id, user_id, user_email, action)
		VALUES ($1, $2, $3, $4)
	`

		_, err = db.Exec(ctx, baseQuery, rowId, user.SubjectID, user.Email, action)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "user": *user, "action": action}).Warn("Failed to log event!")
		}
	}()
}

func GetActionGroup(options *models.GetActionGroupRequest, user *users.User) ([]models.GetActionGroupPart, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	args := make([]any, 0)

	baseQuery := `
		SELECT row_id, user_id, user_email, action, taken_at 
		FROM actions
	`

	if options.Timestamp != nil && options.RowID != nil && *options.RowID != "" {
		args = append(args, time.Unix(*options.Timestamp, 0))
		args = append(args, *options.RowID)
		baseQuery = fmt.Sprintf("%s WHERE (taken_at, row_id) < ($%d, $%d)", baseQuery, len(args)-1, len(args))
	}

	baseQuery = fmt.Sprintf("%s ORDER BY taken_at DESC LIMIT %d", baseQuery, MAX_ACTION_AMOUNT)

	rows, err := db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	actions := make([]models.GetActionGroupPart, 0)
	for rows.Next() {
		var (
			rowId     string
			userId    string
			userEmail string
			action    string
			timestamp time.Time
		)

		if err := rows.Scan(&rowId, &userId, &userEmail, &action, &timestamp); err != nil {
			return nil, false, err
		}

		actions = append(actions, models.GetActionGroupPart{
			RowID:     rowId,
			UserID:    userId,
			UserEmail: userEmail,
			Action:    resolveInserts(action),
			Timestamp: timestamp.Unix(),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	hasMore := (len(actions) > (MAX_ACTION_AMOUNT - 1))
	if hasMore {
		actions = actions[:(MAX_ACTION_AMOUNT)]
	}
	return actions, hasMore, nil
}
