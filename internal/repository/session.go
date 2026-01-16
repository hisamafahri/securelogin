package repository

import (
	"errors"

	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"gorm.io/gorm"
)

var ErrTokenNotFound = errors.New("token not found")

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(
	session *models.Session,
) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) GetByTokenWithRelations(
	token string,
) (*models.Session, error) {
	var session models.Session
	err := r.db.
		Preload("User").
		Preload("AuthenticationRequest").
		Preload("AuthenticationRequest.Application").
		Where("token = ? AND revoked_at IS NULL", token).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) Delete(
	session *models.Session,
) error {
	return r.db.Delete(session).Error
}

func (r *SessionRepository) RevokeByToken(
	token string,
) error {
	result := r.db.Model(&models.Session{}).
		Where("token = ? AND revoked_at IS NULL", token).
		Update("revoked_at", "NOW()")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTokenNotFound
	}
	return nil
}
