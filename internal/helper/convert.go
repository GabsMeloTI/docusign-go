package helper

import (
	"strconv"
)

func ParseStringToInt64(strUserID string) (int64, error) {
	if strUserID == "" {
		return 0, nil
	}
	userID, err := strconv.ParseInt(strUserID, 10, 64)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func ParseStringToInt32(strUserID string) (int32, error) {
	if strUserID == "" {
		return 0, nil
	}
	userID, err := strconv.ParseInt(strUserID, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(userID), nil
}
