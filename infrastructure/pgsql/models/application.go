package models

import (
	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/types"
)

type Application struct {
	ID           uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name         string            `gorm:"type:varchar(255);not null"`
	RedirectURIs types.StringArray `gorm:"type:text[];not null"`
	OriginURIs   types.StringArray `gorm:"type:text[];not null"`
	ClientID     string            `gorm:"type:varchar(255);not null;unique;index"`
	ClientSecret string            `gorm:"type:varchar(255);not null"`
}

func (Application) TableName() string {
	return "applications"
}
