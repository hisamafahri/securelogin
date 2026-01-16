package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"github.com/hisamafahri/securelogin/internal/repository"
	"github.com/hisamafahri/securelogin/pkg/utils"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrSessionRevoked       = errors.New("session has been revoked")
)

type SessionService struct {
	repo *repository.SessionRepository
}

func NewSessionService(
	repo *repository.SessionRepository,
) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) CreateSession(
	authRequestID uuid.UUID,
	userID uuid.UUID,
) (*models.Session, error) {
	token, err := s.generateSecureToken()
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		Token:                   token,
		AuthenticationRequestID: authRequestID,
		UserID:                  userID,
		CreatedAt:               time.Now(),
		ExpiresAt:               time.Now().Add(1 * time.Hour),
	}

	err = s.repo.Create(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) generateSecureToken() (string, error) {
	return utils.GenerateRandomID("session")
}

func (s *SessionService) GetByTokenWithRelations(
	token string,
) (*models.Session, error) {
	session, err := s.repo.GetByTokenWithRelations(token)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionService) ValidateRefreshToken(
	refreshToken string,
) (*models.Session, error) {
	session, err := s.GetByTokenWithRelations(refreshToken)
	if err != nil {
		return nil, ErrRefreshTokenNotFound
	}

	if session.RevokedAt != nil {
		return nil, ErrSessionRevoked
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, ErrRefreshTokenExpired
	}

	return session, nil
}

func (s *SessionService) RotateSession(
	oldSession *models.Session,
) (*models.Session, error) {
	err := s.repo.Delete(oldSession)
	if err != nil {
		return nil, err
	}

	newSession, err := s.CreateSession(
		oldSession.AuthenticationRequestID,
		oldSession.UserID,
	)
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

func (s *SessionService) RevokeToken(
	refreshToken string,
) error {
	err := s.repo.RevokeByToken(refreshToken)
	if err != nil {
		return err
	}
	return nil
}
