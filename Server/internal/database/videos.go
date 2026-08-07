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

	logging.Logger.Info(options.Filters.FilterElement)

	var sortingElement string
	switch options.Filters.FilterElement {
	case models.Title:
		sortingElement = "custom_title"
	case models.Filename:
		sortingElement = "filename"
	case models.DateCreated:
		sortingElement = "file_mod_date"
	case models.DateUploaded:
		sortingElement = "uploaded_at"
	default:
		return nil, false, fmt.Errorf("failed to find correct filter %s", options.Filters.FilterElement)
	}

	baseQuery := fmt.Sprintf(`
		SELECT row_id, s3_id, custom_title, custom_description, visibility, user_id, filename, %s FROM videos WHERE status = 'Complete'
	`, sortingElement)

	if len(options.Filters.FilterTags) > 0 {
		args = append(args, options.Filters.FilterTags)
		baseQuery = fmt.Sprintf(`
        SELECT DISTINCT
            v.row_id,
			v.s3_id,
			v.custom_title,
			v.custom_description,
			v.visibility,
			v.user_id,
			v.filename,
			v.%s
        FROM videos v
        JOIN video_tags vt ON vt.video_id = v.row_id
        WHERE vt.tag_id = ANY($%d)
		AND status = 'Complete'	
		`, sortingElement, len(args))
		if (user == nil) || !user.CanViewPrivateVideos() {
			baseQuery = fmt.Sprintf("%s AND v.visibility = 'Public'", baseQuery)
		}

		if options.Filters.Title != nil {
			args = append(args, "%"+*options.Filters.Title+"%")
			baseQuery = fmt.Sprintf("%s AND v.custom_title ILIKE $%d", baseQuery, len(args))
		}
	} else {
		if (user == nil) || !user.CanViewPrivateVideos() {
			baseQuery = fmt.Sprintf("%s AND visibility = 'Public'", baseQuery)
		}
		if options.Filters.Title != nil {
			args = append(args, "%"+*options.Filters.Title+"%")
			baseQuery = fmt.Sprintf("%s AND custom_title ILIKE $%d", baseQuery, len(args))
		}
	}

	cmp := "<"
	keyWord := "DESC"
	if options.Filters.FilterDirection == models.Ascending {
		cmp = ">"
		keyWord = "ASC"
	}

	if options.Timestamp != nil && options.RowID != nil && *options.RowID != "" {
		switch v := (*options.Timestamp).(type) {
		case int64:
			args = append(args, time.Unix(v, 0))
		case float64:
			args = append(args, time.Unix(int64(v), 0))
		case string:
			args = append(args, v)
		default:
			return nil, false, fmt.Errorf("unsupported timestamp type %T", v)
		}

		args = append(args, options.RowID)
		baseQuery = fmt.Sprintf("%s AND (%s, row_id) %s ($%d, $%d)", baseQuery, sortingElement, cmp, len(args)-1, len(args))
	}

	baseQuery = fmt.Sprintf("%s ORDER BY %s %s, row_id %s", baseQuery, sortingElement, keyWord, keyWord)

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
			visibility  string
			userId      string
			filename    string
			timestamp   any
		)

		if err := rows.Scan(&rowId, &s3Id, &customTitle, &customDesc, &visibility, &userId, &filename, &timestamp); err != nil {
			return nil, false, err
		}

		if t, ok := (timestamp).(time.Time); ok {
			timestamp = t.Unix()
		}

		videos = append(videos, models.GetVideoGroupPart{
			RowID:             rowId,
			S3ID:              s3Id,
			CustomTitle:       customTitle,
			CustomDescription: customDesc,
			Visibility:        visibility,
			UserId:            userId,
			Filename:          filename,
			Timestamp:         timestamp,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	hasMore := (len(videos) > (MAX_ROW_AMOUNT * 3))
	if hasMore {
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

func ChangeVideoVisibility(rowId string, visibility string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	switch visibility {
	case "Public", "Private":
	default:
		return fmt.Errorf("invalid visibility request %s", visibility)
	}

	if _, err := db.Exec(ctx, `
		UPDATE videos
		SET visibility = $1
		WHERE row_id = $2
	`, visibility, rowId); err != nil {
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
		visibility  string
		userId      string
		filename    string
		timestamp   time.Time
	)

	if err := db.QueryRow(ctx, `
		SELECT s3_id, custom_title, custom_description, visibility, user_id, filename, file_mod_date FROM videos
		WHERE row_id = $1
		LIMIT 1 
	`, rowId).Scan(&s3Id, &customTitle, &customDesc, &visibility, &userId, &filename, &timestamp); err != nil {
		return nil, err
	}

	data := &models.GetVideoDataDatabase{
		RowID:             rowId,
		S3ID:              s3Id,
		CustomTitle:       customTitle,
		CustomDescription: customDesc,
		Visibility:        visibility,
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
