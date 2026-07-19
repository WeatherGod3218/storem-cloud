package database

import (
	"context"
	"errors"
	"time"

	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/WeatherGod3218/weather-reels-server/internal/tags"
	"github.com/jackc/pgx/v5/pgconn"
)

const ALREADY_EXISTS_ERROR = "23505"

func GetAllTags() ([]*models.Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	rows, err := db.Query(ctx, `
		SELECT
			row_id,
			name,
			color,
			created_by
		FROM tags
	`)
	if err != nil {
		return nil, err
	}

	var tags []*models.Tag = make([]*models.Tag, 0)

	for rows.Next() {
		var (
			rowId     string
			name      string
			color     string
			createdBy string
		)

		rows.Scan(&rowId, &name, &color, &createdBy)

		tag := &models.Tag{
			TagID:     rowId,
			Name:      name,
			Color:     color,
			CreatedBy: createdBy,
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func CreateTagRow(tagName string, user string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	rowId, err := GenerateUUID()
	if err != nil {
		return false, err
	}

	_, err = db.Exec(ctx, `
		INSERT INTO tags (row_id, name, color, created_by)
		VALUES ($1, $2, $3, $4)
	`, rowId, tagName, tags.GenerateRandomTagHexColor(), "The Wizard!")

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == ALREADY_EXISTS_ERROR {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
