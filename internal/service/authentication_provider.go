package service

import (
	"fmt"
	"net/url"
	"os"

	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"github.com/hisamafahri/securelogin/internal/repository"
)

type ProviderInfo struct {
	Provider models.ProviderType
	ID       uuid.UUID
	Name     string
	Icon     string
}

type AuthenticationProviderService struct {
	repo *repository.AuthenticationProviderRepository
}

func NewAuthenticationProviderService(
	repo *repository.AuthenticationProviderRepository,
) *AuthenticationProviderService {
	return &AuthenticationProviderService{repo: repo}
}

func (s *AuthenticationProviderService) GetProviders(
	applicationID uuid.UUID,
) ([]ProviderInfo, error) {
	providers, err := s.repo.GetByApplicationID(applicationID)
	if err != nil {
		return nil, err
	}

	var providerInfos []ProviderInfo
	for _, provider := range providers {
		providerInfos = append(providerInfos, ProviderInfo{
			Provider: provider.Provider,
			Name:     s.getProviderDisplayName(provider.Provider),
			Icon:     s.getProviderIcon(provider.Provider),
			ID:       provider.ID,
		})
	}

	return providerInfos, nil
}

func (s *AuthenticationProviderService) getProviderDisplayName(
	provider models.ProviderType,
) string {
	switch provider {
	case models.ProviderGoogle:
		return "Continue with Google"
	case models.ProviderGithub:
		return "Continue with GitHub"
	case models.ProviderMicrosoft:
		return "Continue with Microsoft"
	default:
		return "Continue"
	}
}

func (s *AuthenticationProviderService) getProviderIcon(
	provider models.ProviderType,
) string {
	switch provider {
	case models.ProviderGoogle:
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/><path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/><path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/><path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/></svg>`
	case models.ProviderGithub:
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#181717" d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.17 6.839 9.49.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.603-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.463-1.11-1.463-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0112 6.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.138 20.167 22 16.418 22 12c0-5.523-4.477-10-10-10z"/></svg>`
	case models.ProviderMicrosoft:
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24"><path fill="#F25022" d="M12 2H3v9h9V2Z" /><path fill="#00A4EF" d="M12 12H3v9h9v-9Z" /><path fill="#7FBA00" d="M22 2h-9v9h9V2Z" /><path fill="#FFB900" d="M22 12h-9v9h9v-9Z" /></svg>`
	default:
		return ""
	}
}

func (s *AuthenticationProviderService) GetByID(
	id uuid.UUID,
) (*models.AuthenticationProvider, error) {
	return s.repo.GetByID(id)
}

func (s *AuthenticationProviderService) VerifyProviderForRequest(
	providerID uuid.UUID,
	requestApplicationID uuid.UUID,
) (*models.AuthenticationProvider, error) {
	provider, err := s.repo.GetByID(providerID)
	if err != nil {
		return nil, err
	}

	if provider.ApplicationID != requestApplicationID {
		return nil, fmt.Errorf("provider does not belong to the application")
	}

	return provider, nil
}

func (s *AuthenticationProviderService) BuildOAuthURL(
	provider *models.AuthenticationProvider,
	requestID uuid.UUID,
) (string, error) {
	baseURL := os.Getenv("SYSTEM_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	var authURL string
	params := url.Values{}
	params.Set("client_id", provider.ClientID)
	params.Set("state", requestID.String())
	params.Set("response_type", "code")

	switch provider.Provider {
	case models.ProviderGoogle:
		authURL = "https://accounts.google.com/o/oauth2/v2/auth"
		params.Set("scope", "openid email profile")
		params.Set("redirect_uri", fmt.Sprintf("%s/callback/google", baseURL))
	case models.ProviderGithub:
		authURL = "https://github.com/login/oauth/authorize"
		params.Set("scope", "read:user user:email")
		params.Set("redirect_uri", fmt.Sprintf("%s/callback/github", baseURL))
	case models.ProviderMicrosoft:
		authURL = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
		params.Set("scope", "openid email profile")
		params.Set("redirect_uri", fmt.Sprintf("%s/callback/microsoft", baseURL))
	default:
		return "", fmt.Errorf("unsupported provider type: %s", provider.Provider)
	}

	return fmt.Sprintf("%s?%s", authURL, params.Encode()), nil
}
