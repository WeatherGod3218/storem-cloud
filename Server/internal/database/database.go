package database

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func InitDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addr := os.Getenv("POSTGRES_ADDR")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	connURL := &url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(user, password),
		Host:   addr,
		Path:   "/" + user,
	}

	newDB, err := pgxpool.New(
		ctx,
		connURL.String(),
	)

	if err != nil {
		return err
	}

	err = newDB.Ping(ctx)
	if err != nil {
		return err
	}
	db = newDB

	_, err = db.Exec(ctx, `
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
			s3_id				TEXT NOT NULL,
			user_id 			TEXT NOT NULL,

			status 				video_statuses NOT NULL DEFAULT 'Complete',

			filename			TEXT NOT NULL,
			file_size  			INT NOT NULL,
			file_length			REAL NOT NULL,
			file_mod_date		TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			
			custom_title		TEXT,
			custom_description	TEXT,
			
			last_verified		TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			uploaded_at 		TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS tags (
			row_id 		UUID PRIMARY KEY,
			name 		TEXT NOT NULL,
			created_by	TEXT NOT NULL,
			created_at 	TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS video_tags (
			video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
			tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
			added_by TEXT NOT NULL,
			added_at TIMESTAMP DEFAULT now(),
			PRIMARY KEY (video_id, tag_id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_video_file_mod_date_lookup 
			ON videos(file_mod_date);
	`)

	logging.Logger.Info("S3 Connection has been connected!")

	return nil
}
