package model

import (
	"github.com/google/uuid"
	"time"
)

type PayloadDTO struct {
	ID           uuid.UUID `json:"id"`
	UserID       string    `json:"user_id"`
	UserNickname string    `json:"user_nickname"`
	ExpiryAt     time.Time `json:"expiry_at"`
	AccessKey    int64     `json:"access_key"`
	AccessID     int64     `json:"access_id"`
	TenantID     uuid.UUID `json:"tenant_id"`
}

type Payload struct {
	Who      string    `json:"who"`
	WhoID    string    `json:"who_id"`
	TenantID uuid.UUID `json:"tenant_id"`
	AccessID int64     `json:"access_id"`
}
