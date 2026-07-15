package database

import (
	"context"
	"time"
)

func CreateTagRow(tagName string, user string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	rowId, err := GenerateUUID()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
		INSERT INTO tags (row_id, name, created_by)
		VALUES ($1, $2, $3)
	`, rowId)

	if err != nil {
		return err
	}

	return nil
}
