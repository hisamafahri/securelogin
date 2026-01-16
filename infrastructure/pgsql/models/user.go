package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID              `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ApplicationID  uuid.UUID              `gorm:"type:uuid;not null;uniqueIndex:idx_app_provider_user"`
	Application    Application            `gorm:"foreignKey:ApplicationID;references:ID"`
	ProviderID     uuid.UUID              `gorm:"type:uuid;not null"`
	Provider       AuthenticationProvider `gorm:"foreignKey:ProviderID;references:ID"`
	ProviderUserID string                 `gorm:"type:varchar(255);not null;index;uniqueIndex:idx_app_provider_user"`
	Email          string                 `gorm:"type:varchar(255);not null;index"`
	Name           *string                `gorm:"type:varchar(255)"`
	AvatarURL      *string                `gorm:"type:text"`
	CreatedAt      time.Time              `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time              `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (User) TableName() string {
	return "users"
}
