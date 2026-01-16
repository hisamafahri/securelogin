package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID                      uuid.UUID             `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Token                   string                `gorm:"type:varchar(255);not null;unique;index"`
	AuthenticationRequestID uuid.UUID             `gorm:"type:uuid;not null"`
	AuthenticationRequest   AuthenticationRequest `gorm:"foreignKey:AuthenticationRequestID;references:ID"`
	UserID                  uuid.UUID             `gorm:"type:uuid;not null"`
	User                    User                  `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt               time.Time             `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	RevokedAt               *time.Time            `gorm:"type:timestamp"`
	ExpiresAt               time.Time             `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP + INTERVAL '1 hour'"`
}

func (Session) TableName() string {
	return "sessions"
}
