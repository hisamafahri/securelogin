package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"github.com/hisamafahri/securelogin/internal/repository"
)

type AuthenticationRequestService struct {
	authReqRepo *repository.AuthenticationRequestRepository
}

func NewAuthenticationRequestService(
	authReqRepo *repository.AuthenticationRequestRepository,
) *AuthenticationRequestService {
	return &AuthenticationRequestService{
		authReqRepo: authReqRepo,
	}
}

func (s *AuthenticationRequestService) Create(
	applicationID uuid.UUID,
	responseType string,
	redirectURI string,
	state *string,
	codeChallenge *string,
	codeChallengeMethod *string,
	scopes []string,
) (*models.AuthenticationRequest, error) {
	req := &models.AuthenticationRequest{
		ApplicationID:       applicationID,
		ResponseType:        responseType,
		RedirectURI:         redirectURI,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		State:               state,
		CreatedAt:           time.Now(),
		ExpiresAt:           time.Now().Add(30 * time.Minute),
		Scopes:              scopes,
	}

	err := s.authReqRepo.Create(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *AuthenticationRequestService) GetByID(
	id uuid.UUID,
) (*models.AuthenticationRequest, error) {
	return s.authReqRepo.GetByID(id)
}

func (s *AuthenticationRequestService) UpdateProviderID(
	id uuid.UUID,
	providerID uuid.UUID,
) error {
	return s.authReqRepo.UpdateProviderID(id, providerID)
}

func (s *AuthenticationRequestService) MarkAsCompleted(
	id uuid.UUID,
) error {
	return s.authReqRepo.MarkAsCompleted(id)
}
