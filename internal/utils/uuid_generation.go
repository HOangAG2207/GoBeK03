package utils

import "github.com/google/uuid"

func UuidGenerator() string {
	return uuid.New().String()
}
