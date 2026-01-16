package models

import "github.com/google/uuid"

type ProviderType string

const (
	ProviderGoogle    ProviderType = "google"
	ProviderGithub    ProviderType = "github"
	ProviderMicrosoft ProviderType = "microsoft"
)

type AuthenticationProvider struct {
	ID            uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ApplicationID uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_client_application"`
	Application   Application  `gorm:"foreignKey:ApplicationID;references:ID"`
	Provider      ProviderType `gorm:"type:varchar(255);not null;index"`
	// NOTE: let's assume all providers is an OAuth2 provider for now
	ClientID     string `gorm:"type:varchar(255);not null;uniqueIndex:idx_client_application"`
	ClientSecret string `gorm:"type:varchar(255);not null"`
}

func (AuthenticationProvider) TableName() string {
	return "authentication_providers"
}
