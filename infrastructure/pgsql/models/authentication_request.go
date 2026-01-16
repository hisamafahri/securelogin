package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/types"
)

type AuthenticationRequest struct {
	ID                  uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ApplicationID       uuid.UUID         `gorm:"type:uuid;not null"`
	Application         Application       `gorm:"foreignKey:ApplicationID;references:ID"`
	ProviderID          *uuid.UUID        `gorm:"type:uuid"`
	ResponseType        string            `gorm:"type:varchar(255);not null"`
	RedirectURI         string            `gorm:"type:text;not null"`
	Scopes              types.StringArray `gorm:"type:text[];not null"`
	State               *string           `gorm:"type:varchar(255)"`
	CodeChallenge       *string           `gorm:"type:text"`
	CodeChallengeMethod *string           `gorm:"type:varchar(255)"`
	CreatedAt           time.Time         `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ExpiresAt           time.Time         `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP + INTERVAL '30 minutes'"`
	CompletedAt         *time.Time        `gorm:"type:timestamp"`
}

func (AuthenticationRequest) TableName() string {
	return "authentication_requests"
}
