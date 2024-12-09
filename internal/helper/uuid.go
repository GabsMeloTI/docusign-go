package helper

import "github.com/google/uuid"

func GetUUidStr() string {
	return uuid.New().String()
}

func GetUUID() uuid.UUID {
	return uuid.New()
}
