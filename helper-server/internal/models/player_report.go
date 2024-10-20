package models

import "github.com/google/uuid"

type PlayerReport struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;not null;unique;default:gen_random_uuid()"`

	SourcePlayerRefer uuid.UUID `gorm:"not null"`
	SourcePlayer      *Player   `gorm:"foreignKey:SourcePlayerRefer"`

	TargetPlayerRefer uuid.UUID `gorm:"not null"`
	TargetPlayer      *Player   `gorm:"foreignKey:TargetPlayerRefer"`

	Subject string
	Content string
}
