package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthorizationCode struct {
	ID                      uuid.UUID             `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code                    string                `gorm:"type:varchar(255);not null;unique;index"`
	AuthenticationRequestID uuid.UUID             `gorm:"type:uuid;not null"`
	AuthenticationRequest   AuthenticationRequest `gorm:"foreignKey:AuthenticationRequestID;references:ID"`
	UserID                  uuid.UUID             `gorm:"type:uuid;not null"`
	User                    User                  `gorm:"foreignKey:UserID;references:ID"`
	UsedAt                  *time.Time            `gorm:"type:timestamp;"`
	CreatedAt               time.Time             `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ExpiresAt               time.Time             `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP + INTERVAL '10 minutes'"`
}

func (AuthorizationCode) TableName() string {
	return "authorization_codes"
}
