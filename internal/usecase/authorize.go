package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"github.com/hisamafahri/securelogin/internal/service"
)

var (
	ErrClientNotFound     = errors.New("application is not found")
	ErrInvalidRedirectURI = errors.New(
		"redirect_uri not registered for this application",
	)
)

type AuthorizeUsecase struct {
	appService     *service.ApplicationService
	authReqService *service.AuthenticationRequestService
}

func NewAuthorizeUsecase(
	appService *service.ApplicationService,
	authReqService *service.AuthenticationRequestService,
) *AuthorizeUsecase {
	return &AuthorizeUsecase{
		appService:     appService,
		authReqService: authReqService,
	}
}

func (u *AuthorizeUsecase) ValidateAuthorizeRequest(
	clientID string,
	redirectURI string,
) error {
	app, err := u.appService.GetByClientID(clientID)
	if err != nil {
		return ErrClientNotFound
	}

	if !u.appService.IsRedirectURIValid(app, redirectURI) {
		return ErrInvalidRedirectURI
	}

	return nil
}

func (u *AuthorizeUsecase) ValidateAndGetApplication(
	clientID string,
	redirectURI string,
) (*models.Application, error) {
	app, err := u.appService.GetByClientID(clientID)
	if err != nil {
		return nil, ErrClientNotFound
	}

	if !u.appService.IsRedirectURIValid(app, redirectURI) {
		return nil, ErrInvalidRedirectURI
	}

	return app, nil
}

func (u *AuthorizeUsecase) CreateAuthenticationRequest(
	applicationID uuid.UUID,
	responseType string,
	redirectURI string,
	state *string,
	codeChallenge *string,
	codeChallengeMethod *string,
	scopes []string,
) (*models.AuthenticationRequest, error) {
	return u.authReqService.Create(
		applicationID,
		responseType,
		redirectURI,
		state,
		codeChallenge,
		codeChallengeMethod,
		scopes,
	)
}
