package database

import (
	"context"
	"time"

	"github.com/WeatherGod3218/weather-reels-server/internal/models"
)

func AddTagToVideo(videoId string, tagId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := db.Exec(ctx, `
		INSERT INTO video_tags (video_id, tag_id, added_by)
		VALUES ($1, $2, $3)
	`, videoId, tagId, "No Oauth Yet :("); err != nil {
		return err
	}

	return nil
}

func RemoveTagFromVideo(videoId string, tagId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := db.Exec(ctx, `
		DELETE FROM video_tags
		WHERE video_id = $1
		AND tag_id = $2
	`, videoId, tagId); err != nil {
		return err
	}

	return nil
}

func GetAllTagsOnVideo(videoId string) ([]*models.VideoTag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := db.Query(ctx, `
		SELECT
			t.row_id,
			t.name,
			t.color
		FROM videos v
		JOIN video_tags vt ON vt.video_id = v.row_id
		JOIN tags t ON t.row_id = vt.tag_id
		WHERE v.row_id = $1
	`, videoId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.VideoTag = make([]*models.VideoTag, 0)

	for rows.Next() {
		var (
			id    string
			name  string
			color string
		)

		if err := rows.Scan(&id, &name, &color); err != nil {
			return nil, err
		}

		tags = append(tags, &models.VideoTag{
			TagID: id,
			Name:  name,
			Color: color,
		})
	}
	return tags, nil
}
