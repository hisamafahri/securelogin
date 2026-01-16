package usecase

import (
	"errors"
	"slices"
	"strings"

	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"github.com/hisamafahri/securelogin/internal/service"
	"github.com/hisamafahri/securelogin/pkg/utils"
)

var (
	ErrInvalidGrantType      = errors.New("invalid grant type")
	ErrRedirectURIMismatch   = errors.New("redirect uri mismatch")
	ErrAuthorizationCodeUsed = errors.New("authorization code already used")
	ErrCodeVerifierRequired  = errors.New("code verifier required")
	ErrInvalidCodeVerifier   = errors.New("invalid code verifier")
	ErrInvalidScope          = errors.New(
		"requested scope is invalid or exceeds originally granted scope",
	)
)

type TokenExchangeUsecase struct {
	appService      *service.ApplicationService
	authCodeService *service.AuthorizationCodeService
	sessionService  *service.SessionService
}

func NewTokenExchangeUsecase(
	appService *service.ApplicationService,
	authCodeService *service.AuthorizationCodeService,
	sessionService *service.SessionService,
) *TokenExchangeUsecase {
	return &TokenExchangeUsecase{
		appService:      appService,
		authCodeService: authCodeService,
		sessionService:  sessionService,
	}
}

func (u *TokenExchangeUsecase) ExchangeCodeForToken(
	grantType string,
	clientID string,
	clientSecret string,
	code string,
	redirectURI string,
	codeVerifier string,
) (*models.Session, *models.User, *models.AuthenticationRequest, error) {
	if grantType != "authorization_code" {
		return nil, nil, nil, ErrInvalidGrantType
	}

	_, err := u.appService.ValidateClientCredentials(clientID, clientSecret)
	if err != nil {
		return nil, nil, nil, err
	}

	authCode, err := u.authCodeService.GetByCodeWithRelations(code)
	if err != nil {
		return nil, nil, nil, err
	}

	if authCode.AuthenticationRequest.RedirectURI != "" && redirectURI != "" {
		if authCode.AuthenticationRequest.RedirectURI != redirectURI {
			return nil, nil, nil, ErrRedirectURIMismatch
		}
	}

	if authCode.AuthenticationRequest.CodeChallenge != nil {
		if codeVerifier == "" {
			return nil, nil, nil, ErrCodeVerifierRequired
		}

		codeChallengeMethod := "plain"
		if authCode.AuthenticationRequest.CodeChallengeMethod != nil {
			codeChallengeMethod = *authCode.AuthenticationRequest.CodeChallengeMethod
		}

		err := utils.VerifyPKCE(
			*authCode.AuthenticationRequest.CodeChallenge,
			codeChallengeMethod,
			codeVerifier,
		)
		if err != nil {
			return nil, nil, nil, ErrInvalidCodeVerifier
		}
	}

	err = u.authCodeService.ValidateAndMarkAsUsed(authCode)
	if err != nil {
		if errors.Is(err, service.ErrAuthorizationCodeUsed) {
			return nil, nil, nil, ErrAuthorizationCodeUsed
		}
		return nil, nil, nil, err
	}

	session, err := u.sessionService.CreateSession(
		authCode.AuthenticationRequestID,
		authCode.UserID,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	return session, &authCode.User, &authCode.AuthenticationRequest, nil
}

func (u *TokenExchangeUsecase) RefreshAccessToken(
	grantType string,
	clientID string,
	clientSecret string,
	refreshToken string,
	requestedScope string,
) (*models.Session, *models.User, *models.AuthenticationRequest, []string, error) {
	if grantType != "refresh_token" {
		return nil, nil, nil, nil, ErrInvalidGrantType
	}

	_, err := u.appService.ValidateClientCredentials(clientID, clientSecret)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	oldSession, err := u.sessionService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if oldSession.AuthenticationRequest.Application.ClientID != clientID {
		return nil, nil, nil, nil, service.ErrInvalidClientCredentials
	}

	originalScopes := oldSession.AuthenticationRequest.Scopes
	finalScopes := originalScopes

	if requestedScope != "" {
		requestedScopes := strings.Split(requestedScope, " ")

		for _, reqScope := range requestedScopes {
			if !slices.Contains(originalScopes, reqScope) {
				return nil, nil, nil, nil, ErrInvalidScope
			}
		}

		finalScopes = requestedScopes
	}

	newSession, err := u.sessionService.RotateSession(oldSession)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return newSession, &oldSession.User, &oldSession.AuthenticationRequest, finalScopes, nil
}
