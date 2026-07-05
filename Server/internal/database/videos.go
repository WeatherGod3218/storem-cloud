package database

import (
	"context"
	"time"
)

func VerifyVideoList(videos []string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	rows, err := db.Query(ctx, `
		SELECT filename
			FROM videos
			WHERE filename = ALL($1);
	`, videos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notBackedUp []string

	found := make(map[string]struct{}, len(videos))
	for rows.Next() {
		var file string
		if err := rows.Scan(&file); err != nil {
			return nil, err
		}
		notBackedUp = append(notBackedUp, file)
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

func InitVideos() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := db.Exec(ctx, `
		DO $$ BEGIN
			CREATE TYPE video_statuses AS ENUM ('Pending', 'Complete', 'Failed', 'Removed');
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

			filename			TEXT NOT NULL,
			status 				video_statuses NOT NULL DEFAULT 'Pending',						
			filesize  			INT NOT NULL,
			video_s3_url		TEXT NOT NULL,
			thumbnail_s3_url	TEXT NOT NULL,

			custom_title	TEXT,

			last_verified	TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			uploaded_at 	TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
