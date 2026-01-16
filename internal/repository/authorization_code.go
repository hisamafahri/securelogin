package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"gorm.io/gorm"
)

type AuthorizationCodeRepository struct {
	db *gorm.DB
}

func NewAuthorizationCodeRepository(db *gorm.DB) *AuthorizationCodeRepository {
	return &AuthorizationCodeRepository{db: db}
}

func (r *AuthorizationCodeRepository) Create(
	authCode *models.AuthorizationCode,
) error {
	return r.db.Create(authCode).Error
}

func (r *AuthorizationCodeRepository) GetByCode(
	code string,
) (*models.AuthorizationCode, error) {
	var authCode models.AuthorizationCode
	err := r.db.Where("code = ?", code).First(&authCode).Error
	if err != nil {
		return nil, err
	}
	return &authCode, nil
}

func (r *AuthorizationCodeRepository) GetByCodeWithRelations(
	code string,
) (*models.AuthorizationCode, error) {
	var authCode models.AuthorizationCode
	err := r.db.
		Preload("AuthenticationRequest").
		Preload("User").
		Where("code = ?", code).
		First(&authCode).Error
	if err != nil {
		return nil, err
	}
	return &authCode, nil
}

func (r *AuthorizationCodeRepository) MarkAsUsed(
	id uuid.UUID,
) error {
	now := time.Now()
	return r.db.Model(&models.AuthorizationCode{}).
		Where("id = ?", id).
		Update("used_at", now).Error
}
