package repository

import (
	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"gorm.io/gorm"
)

type AuthenticationProviderRepository struct {
	db *gorm.DB
}

func NewAuthenticationProviderRepository(
	db *gorm.DB,
) *AuthenticationProviderRepository {
	return &AuthenticationProviderRepository{db: db}
}

func (r *AuthenticationProviderRepository) GetByApplicationID(
	applicationID uuid.UUID,
) ([]models.AuthenticationProvider, error) {
	var providers []models.AuthenticationProvider
	err := r.db.Where("application_id = ?", applicationID).Find(&providers).Error
	if err != nil {
		return nil, err
	}
	return providers, nil
}

func (r *AuthenticationProviderRepository) GetByID(
	id uuid.UUID,
) (*models.AuthenticationProvider, error) {
	var provider models.AuthenticationProvider
	err := r.db.Where("id = ?", id).First(&provider).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}
