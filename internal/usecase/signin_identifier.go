package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/internal/service"
)

type SigninIdentifierUsecase struct {
	authReqService  *service.AuthenticationRequestService
	providerService *service.AuthenticationProviderService
}

func NewSigninIdentifierUsecase(
	authReqService *service.AuthenticationRequestService,
	providerService *service.AuthenticationProviderService,
) *SigninIdentifierUsecase {
	return &SigninIdentifierUsecase{
		authReqService:  authReqService,
		providerService: providerService,
	}
}

func (u *SigninIdentifierUsecase) Execute(
	providerID uuid.UUID,
	requestID uuid.UUID,
) (string, error) {
	authReq, err := u.authReqService.GetByID(requestID)
	if err != nil {
		return "", fmt.Errorf("authentication request not found")
	}

	provider, err := u.providerService.VerifyProviderForRequest(
		providerID,
		authReq.ApplicationID,
	)
	if err != nil {
		return "", fmt.Errorf("provider verification failed: %w", err)
	}

	err = u.authReqService.UpdateProviderID(requestID, providerID)
	if err != nil {
		return "", fmt.Errorf("failed to update provider id: %w", err)
	}

	oauthURL, err := u.providerService.BuildOAuthURL(provider, requestID)
	if err != nil {
		return "", fmt.Errorf("failed to build oauth url: %w", err)
	}

	return oauthURL, nil
}
