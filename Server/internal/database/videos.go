package database

import (
	"context"
	"fmt"
	"time"

	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/WeatherGod3218/weather-reels-server/internal/users"
)

const MAX_ROW_AMOUNT int = 3

func VerifyVideoList(videos []string) ([]string, error) {
	if len(videos) == 0 {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	rows, err := db.Query(ctx, `
		SELECT filename
			FROM videos
			WHERE filename = ANY($1) AND status IS DISTINCT FROM 'Failed';
	`, videos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	found := make(map[string]struct{}, len(videos))
	for rows.Next() {
		var file string
		if err := rows.Scan(&file); err != nil {
			return nil, err
		}
		found[file] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var missing []string
	for _, v := range videos {
		if _, ok := found[v]; !ok {
			missing = append(missing, v)
		}
	}

	if len(found) > 0 {
		existing := make([]string, 0, len(found))
		for h := range found {
			existing = append(existing, h)
		}
		_, err := db.Exec(ctx, `
            UPDATE videos
                SET last_verified = NOW()
                WHERE filename = ANY($1);
        `, existing)

		if err != nil {
			return nil, err
		}
	}
	return missing, nil
}

func GetVideoGroup(options *models.GetVideoGroupRequest, user *users.User) ([]models.GetVideoGroupPart, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	args := make([]any, 0)

	baseQuery := `
		SELECT row_id, s3_id, custom_title, custom_description, user_id, filename, file_mod_date FROM videos WHERE status = 'Complete'
	`

	if user == nil || !user.CanViewPrivateVideos() {
		baseQuery = fmt.Sprintf("%s AND visibility = 'Public'", baseQuery)
	}

	cmp := "<"
	if options.OrderAscending {
		cmp = ">"
	}

	if options.Timestamp != nil && options.RowID != nil {
		args = append(args, time.Unix(*options.Timestamp, 0))
		args = append(args, options.RowID)
		baseQuery = fmt.Sprintf("%s AND (file_mod_date, row_id) %s ($%d, $%d)", baseQuery, cmp, len(args)-1, len(args))
	}

	if options.OrderAscending == true {
		baseQuery = fmt.Sprintf("%s ORDER BY file_mod_date ASC, row_id ASC", baseQuery)
	} else {
		baseQuery = fmt.Sprintf("%s ORDER BY file_mod_date DESC, row_id DESC", baseQuery)
	}

	//add limit
	args = append(args, (MAX_ROW_AMOUNT*3)+1)
	baseQuery = fmt.Sprintf("%s LIMIT $%d", baseQuery, len(args))

	rows, err := db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	videos := make([]models.GetVideoGroupPart, 0)
	for rows.Next() {
		var (
			rowId       string
			s3Id        string
			customTitle *string
			customDesc  *string
			userId      string
			filename    string
			timestamp   time.Time
		)

		if err := rows.Scan(&rowId, &s3Id, &customTitle, &customDesc, &userId, &filename, &timestamp); err != nil {
			return nil, false, err
		}

		videos = append(videos, models.GetVideoGroupPart{
			RowID:             rowId,
			S3ID:              s3Id,
			CustomTitle:       customTitle,
			CustomDescription: customDesc,
			UserId:            userId,
			Filename:          filename,
			Timestamp:         timestamp.Unix(),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	hasMore := (len(videos) > (MAX_ROW_AMOUNT * 3))
	if hasMore {
		logging.Logger.Info("does not have more!")
		videos = videos[:(MAX_ROW_AMOUNT * 3)]
	}
	return videos, hasMore, nil
}

func GetRandomVideoData(user *users.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	baseQuery := `
		SELECT row_id FROM videos WHERE status = 'Complete'
	`
	if user == nil || !user.CanViewPrivateVideos() {
		baseQuery = fmt.Sprintf("%s AND visibility = 'Public'", baseQuery)
	}

	baseQuery = fmt.Sprintf("%s ORDER BY RANDOM() LIMIT 1", baseQuery)
	var (
		rowId string
	)

	if err := db.QueryRow(ctx, baseQuery).Scan(&rowId); err != nil {
		return "", err
	}

	return rowId, nil
}

func ChangeVideoTitle(rowId string, title string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := db.Exec(ctx, `
		UPDATE videos
		SET custom_title = $1
		WHERE row_id = $2
	`, title, rowId); err != nil {
		return err
	}

	return nil
}

func ChangeVideoDescription(rowId string, desc string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := db.Exec(ctx, `
		UPDATE videos
		SET custom_description = $1
		WHERE row_id = $2
	`, desc, rowId); err != nil {
		return err
	}

	return nil
}

func GetVideoData(rowId string) (*models.GetVideoDataDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var (
		s3Id        string
		customTitle *string
		customDesc  *string
		userId      string
		filename    string
		timestamp   time.Time
	)

	if err := db.QueryRow(ctx, `
		SELECT s3_id, custom_title, custom_description, user_id, filename, file_mod_date FROM videos
		WHERE row_id = $1
		LIMIT 1 
	`, rowId).Scan(&s3Id, &customTitle, &customDesc, &userId, &filename, &timestamp); err != nil {
		return nil, err
	}

	data := &models.GetVideoDataDatabase{
		RowID:             rowId,
		S3ID:              s3Id,
		CustomTitle:       customTitle,
		CustomDescription: customDesc,
		UserID:            userId,
		Filename:          filename,
		Timestamp:         timestamp.Unix(),
	}

	return data, nil
}

func CreateVideoRow(video models.VideoDatabaseEntry) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	rowId, err := GenerateUUID()
	if err != nil {
		return "", err
	}

	_, err = db.Exec(ctx, `
		INSERT INTO videos (row_id, s3_id, user_id, filename, file_size, file_length, file_mod_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, rowId, video.VideoID, video.UserID, video.Filename, video.FileSize, video.FileLength, video.FileModDate)

	if err != nil {
		return "", err
	}
	return rowId, nil
}
