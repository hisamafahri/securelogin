package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"gorm.io/gorm"
)

type AuthenticationRequestRepository struct {
	db *gorm.DB
}

func NewAuthenticationRequestRepository(
	db *gorm.DB,
) *AuthenticationRequestRepository {
	return &AuthenticationRequestRepository{db: db}
}

func (r *AuthenticationRequestRepository) Create(
	req *models.AuthenticationRequest,
) error {
	return r.db.Create(req).Error
}

func (r *AuthenticationRequestRepository) GetByID(
	id uuid.UUID,
) (*models.AuthenticationRequest, error) {
	var req models.AuthenticationRequest
	err := r.db.Where("id = ?", id).First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *AuthenticationRequestRepository) UpdateProviderID(
	id uuid.UUID,
	providerID uuid.UUID,
) error {
	return r.db.Model(&models.AuthenticationRequest{}).
		Where("id = ?", id).
		Update("provider_id", providerID).Error
}

func (r *AuthenticationRequestRepository) MarkAsCompleted(
	id uuid.UUID,
) error {
	now := time.Now()
	return r.db.Model(&models.AuthenticationRequest{}).
		Where("id = ?", id).
		Update("completed_at", now).Error
}
