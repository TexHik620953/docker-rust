package models

import "github.com/google/uuid"

type Player struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid;not null;unique;default:gen_random_uuid()"`
	SteamID    uint64    `gorm:"not null"`
	PlayerName string    `gorm:""`

	Balance int64 `gorm:"not null;default:0"`
}
