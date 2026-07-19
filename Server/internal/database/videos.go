package database

import (
	"context"
	"time"

	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/jackc/pgx/v5"
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

func GetVideoGroup(offset *time.Time, rowID string) ([]models.GetVideoGroupPart, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var (
		rows pgx.Rows
		err  error
	)
	if offset == nil {
		rows, err = db.Query(ctx, `
			SELECT row_id, s3_id, custom_title, custom_description, user_id, filename, file_mod_date FROM videos
			WHERE status = 'Complete'
			ORDER BY file_mod_date DESC, row_id DESC
			LIMIT $1 
		`, (MAX_ROW_AMOUNT*3)+1)
	} else {
		rows, err = db.Query(ctx, `
			SELECT row_id, s3_id, custom_title, custom_description, user_id, filename, file_mod_date FROM videos
			WHERE status = 'Complete'
			AND (file_mod_date, row_id) < ($1, $2)
			ORDER BY file_mod_date DESC, row_id DESC
			LIMIT $3 
		`, offset, rowID, (MAX_ROW_AMOUNT*3)+1)
	}
	defer rows.Close()

	if err != nil {
		return nil, false, err
	}

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
			S3Id:              s3Id,
			CustomTitle:       customTitle,
			CustomDescription: customDesc,
			UserId:            userId,
			Filename:          filename,
			Timestamp:         timestamp,
		})
	}

	hasMore := (len(videos) > (MAX_ROW_AMOUNT * 3))
	if hasMore {
		videos = videos[:(MAX_ROW_AMOUNT * 3)]
	}
	return videos, hasMore, nil
}

func GetRandomVideoData() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var (
		rowId string
	)

	if err := db.QueryRow(ctx, `
		SELECT row_id FROM videos
		WHERE status = 'Complete'
		ORDER BY RANDOM()
		LIMIT 1
	`).Scan(&rowId); err != nil {
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
		UserId:            userId,
		Filename:          filename,
		Timestamp:         timestamp,
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
	`, rowId, video.VideoId, video.UserId, video.Filename, video.FileSize, video.FileLength, video.FileModDate)

	if err != nil {
		return "", err
	}
	return rowId, nil
}
