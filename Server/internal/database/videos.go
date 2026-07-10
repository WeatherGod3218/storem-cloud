package database

import (
	"context"
	"time"

	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/jackc/pgx/v5"
)

const MAX_GROUP_SIZE = 10

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

func AbortVideoUpload(row_id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := db.Exec(ctx, `
		UPDATE videos
		SET status = 'Failed'
		WHERE row_id = $2 AND status IS DISTINCT FROM 'Complete'
	`, row_id)
	if err != nil {
		return err
	}

	return nil
}

func GetVideoGroup(offset *time.Time, rowID string) ([]models.VideoGroupPart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var (
		rows pgx.Rows
		err  error
	)
	if offset == nil {
		rows, err = db.Query(ctx, `
			SELECT row_id, filename, uploaded_at FROM videos
			WHERE status = 'Complete'
			ORDER BY uploaded_at DESC, row_id DESC
			LIMIT 10 
		`)
	} else {
		rows, err = db.Query(ctx, `
			SELECT row_id, filename, uploaded_at FROM videos
			WHERE status = 'Complete'
			AND (uploaded_at, row_id) < ($1, $2)
			ORDER BY uploaded_at DESC, row_id DESC
			LIMIT 10 
		`, offset, rowID)
	}

	if err != nil {
		return nil, err
	}

	videos := make([]models.VideoGroupPart, 0)
	for rows.Next() {
		var (
			row       string
			filename  string
			timestamp time.Time
		)

		if err := rows.Scan(&row, &filename, &timestamp); err != nil {
			return nil, err
		}

		videos = append(videos, models.VideoGroupPart{
			RowID:     row,
			FileName:  filename,
			Timestamp: timestamp,
		})
	}

	return videos, nil
}

func CompleteVideoUpload(row_id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := db.Exec(ctx, `
		UPDATE videos
		SET status = 'Complete'
		WHERE row_id = $1
	`, row_id)
	if err != nil {
		return err
	}

	return nil
}

func StartVideoUpload(video models.VideoBackupProcessed) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	rowId, err := GenerateUUID()
	if err != nil {
		return "", err
	}

	_, err = db.Exec(ctx, `
		INSERT INTO videos (row_id, filename, filesize, filelength, video_s3_url, thumbnail_s3_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (filename)
		DO UPDATE SET
			status = 'Pending',
			last_verified = NOW();

	`, rowId, video.FileName, video.FileSize, video.FileLength, video.VideoS3URL, video.ThumbnailS3URL)

	if err != nil {
		return "", err
	}
	return rowId, nil
}

func InitVideos() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := db.Exec(ctx, `
		DO $$ BEGIN
			CREATE TYPE video_statuses AS ENUM ('Pending', 'Complete', 'Failed', 'Removed', 'Excluded');
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$;
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS videos (
			row_id 				UUID PRIMARY KEY,
			status 				video_statuses NOT NULL DEFAULT 'Pending',

			filename			TEXT NOT NULL,
			filesize  			INT NOT NULL,
			filelength			REAL NOT NULL,

			video_s3_url		TEXT NOT NULL,
			thumbnail_s3_url	TEXT NOT NULL,

			custom_title	TEXT,

			last_verified	TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			uploaded_at 	TIMESTAMPTZ NOT NULL DEFAULT NOW(),

			UNIQUE(filename)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_video_filename_lookup 
			ON videos(filename);
	`)

	return err
}
