package models

import "github.com/google/uuid"

type FeedbackType int

const (
	FEEDBACK_TYPE_GENERAL   FeedbackType = iota
	FEEDBACK_TYPE_BUG       FeedbackType = iota
	FEEDBACK_TYPE_CHEAT     FeedbackType = iota
	FEEDBACK_TYPE_ABUSE     FeedbackType = iota
	FEEDBACK_TYPE_IDEA      FeedbackType = iota
	FEEDBACK_TYPE_OFFENSIVE FeedbackType = iota
)

type PlayerFeedback struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;not null;unique;default:gen_random_uuid()"`

	PlayerRefer uuid.UUID `gorm:"not null"`
	Player      *Player   `gorm:"foreignKey:PlayerRefer"`

	FeedbackType int `gorm:"not null;default:0"`
	Subject      string
	Content      string
}
