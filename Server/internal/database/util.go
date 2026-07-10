package database

import (
	"github.com/google/uuid"
)

func GenerateUUID() (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", nil
	}

	return uuid.String(), nil
}
