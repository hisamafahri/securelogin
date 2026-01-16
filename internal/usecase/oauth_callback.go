package usecase

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/internal/service"
)

type OAuthCallbackUsecase struct {
	authReqService  *service.AuthenticationRequestService
	providerService *service.AuthenticationProviderService
	oauthService    *service.OAuthExchangeService
	userService     *service.UserService
	authCodeService *service.AuthorizationCodeService
	sessionService  *service.SessionService
}

func NewOAuthCallbackUsecase(
	authReqService *service.AuthenticationRequestService,
	providerService *service.AuthenticationProviderService,
	oauthService *service.OAuthExchangeService,
	userService *service.UserService,
	authCodeService *service.AuthorizationCodeService,
	sessionService *service.SessionService,
) *OAuthCallbackUsecase {
	return &OAuthCallbackUsecase{
		authReqService:  authReqService,
		providerService: providerService,
		oauthService:    oauthService,
		userService:     userService,
		authCodeService: authCodeService,
		sessionService:  sessionService,
	}
}

func (u *OAuthCallbackUsecase) Execute(
	code string,
	state string,
) (string, error) {
	requestID, err := uuid.Parse(state)
	if err != nil {
		return "", fmt.Errorf("invalid state parameter")
	}

	authReq, err := u.authReqService.GetByID(requestID)
	if err != nil {
		return "", fmt.Errorf("authentication request not found")
	}

	if authReq.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("authentication request has expired")
	}

	if authReq.ProviderID == nil {
		return "", fmt.Errorf("provider not set for authentication request")
	}

	provider, err := u.providerService.GetByID(*authReq.ProviderID)
	if err != nil {
		return "", fmt.Errorf("provider not found")
	}

	accessToken, err := u.oauthService.ExchangeCodeForToken(provider, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	userProfile, err := u.oauthService.GetUserProfile(provider, accessToken)
	if err != nil {
		return "", fmt.Errorf("failed to get user profile: %w", err)
	}

	user, err := u.userService.FindOrCreateUser(
		authReq.ApplicationID,
		*authReq.ProviderID,
		userProfile.ProviderUserID,
		userProfile.Email,
		&userProfile.Name,
		&userProfile.AvatarURL,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	err = u.authReqService.MarkAsCompleted(authReq.ID)
	if err != nil {
		return "", fmt.Errorf("failed to mark request as completed: %w", err)
	}

	switch authReq.ResponseType {
	case "token":
		session, err := u.sessionService.CreateSession(
			authReq.ID,
			user.ID,
		)
		if err != nil {
			return "", fmt.Errorf("failed to create session: %w", err)
		}

		redirectURL, err := url.Parse(authReq.RedirectURI)
		if err != nil {
			return "", fmt.Errorf("invalid redirect uri")
		}

		fragment := url.Values{}
		fragment.Set("access_token", session.Token)
		fragment.Set("token_type", "Bearer")
		if authReq.State != nil {
			fragment.Set("state", *authReq.State)
		}
		redirectURL.Fragment = fragment.Encode()

		return redirectURL.String(), nil
	case "code":
		authCode, err := u.authCodeService.CreateAuthorizationCode(
			authReq.ID,
			user.ID,
		)
		if err != nil {
			return "", fmt.Errorf("failed to create authorization code: %w", err)
		}

		redirectURL, err := url.Parse(authReq.RedirectURI)
		if err != nil {
			return "", fmt.Errorf("invalid redirect uri")
		}

		query := redirectURL.Query()
		query.Set("code", authCode)
		if authReq.State != nil {
			query.Set("state", *authReq.State)
		}
		redirectURL.RawQuery = query.Encode()

		return redirectURL.String(), nil
	default:
		return "", fmt.Errorf("unsupported response type")
	}
}
