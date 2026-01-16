package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"github.com/hisamafahri/securelogin/internal/repository"
	"github.com/hisamafahri/securelogin/pkg/utils"
	"gorm.io/gorm"
)

var (
	ErrAuthorizationCodeNotFound = errors.New("authorization code not found")
	ErrAuthorizationCodeExpired  = errors.New("authorization code expired")
	ErrAuthorizationCodeUsed     = errors.New("authorization code already used")
)

type AuthorizationCodeService struct {
	repo *repository.AuthorizationCodeRepository
}

func NewAuthorizationCodeService(
	repo *repository.AuthorizationCodeRepository,
) *AuthorizationCodeService {
	return &AuthorizationCodeService{repo: repo}
}

func (s *AuthorizationCodeService) CreateAuthorizationCode(
	authRequestID uuid.UUID,
	userID uuid.UUID,
) (string, error) {
	code, err := s.generateSecureCode()
	if err != nil {
		return "", err
	}

	authCode := &models.AuthorizationCode{
		Code:                    code,
		AuthenticationRequestID: authRequestID,
		UserID:                  userID,
		UsedAt:                  nil,
		CreatedAt:               time.Now(),
		ExpiresAt:               time.Now().Add(10 * time.Minute),
	}

	err = s.repo.Create(authCode)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (s *AuthorizationCodeService) generateSecureCode() (string, error) {
	return utils.GenerateRandomID("auth")
}

func (s *AuthorizationCodeService) GetByCodeWithRelations(
	code string,
) (*models.AuthorizationCode, error) {
	authCode, err := s.repo.GetByCodeWithRelations(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuthorizationCodeNotFound
		}
		return nil, err
	}
	return authCode, nil
}

func (s *AuthorizationCodeService) ValidateAndMarkAsUsed(
	authCode *models.AuthorizationCode,
) error {
	if authCode.UsedAt != nil {
		return ErrAuthorizationCodeUsed
	}

	if time.Now().After(authCode.ExpiresAt) {
		return ErrAuthorizationCodeExpired
	}

	return s.repo.MarkAsUsed(authCode.ID)
}
