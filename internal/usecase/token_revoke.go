package usecase

import (
	"github.com/hisamafahri/securelogin/internal/service"
)

type TokenRevokeUsecase struct {
	appService     *service.ApplicationService
	sessionService *service.SessionService
}

func NewTokenRevokeUsecase(
	appService *service.ApplicationService,
	sessionService *service.SessionService,
) *TokenRevokeUsecase {
	return &TokenRevokeUsecase{
		appService:     appService,
		sessionService: sessionService,
	}
}

func (u *TokenRevokeUsecase) RevokeToken(
	clientID string,
	clientSecret string,
	refreshToken string,
) error {
	_, err := u.appService.ValidateClientCredentials(clientID, clientSecret)
	if err != nil {
		return err
	}

	err = u.sessionService.RevokeToken(refreshToken)
	if err != nil {
		return err
	}

	return nil
}
