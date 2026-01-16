package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
)

type UserProfile struct {
	ProviderUserID string
	Email          string
	Name           string
	AvatarURL      string
}

type OAuthExchangeService struct{}

func NewOAuthExchangeService() *OAuthExchangeService {
	return &OAuthExchangeService{}
}

func (s *OAuthExchangeService) ExchangeCodeForToken(
	provider *models.AuthenticationProvider,
	code string,
) (string, error) {
	baseURL := os.Getenv("SYSTEM_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	var tokenURL string
	var redirectURI string

	switch provider.Provider {
	case models.ProviderGoogle:
		tokenURL = "https://oauth2.googleapis.com/token"
		redirectURI = fmt.Sprintf("%s/callback/google", baseURL)
	case models.ProviderGithub:
		tokenURL = "https://github.com/login/oauth/access_token"
		redirectURI = fmt.Sprintf("%s/callback/github", baseURL)
	case models.ProviderMicrosoft:
		tokenURL = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
		redirectURI = fmt.Sprintf("%s/callback/microsoft", baseURL)
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider.Provider)
	}

	data := url.Values{}
	data.Set("client_id", provider.ClientID)
	data.Set("client_secret", provider.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest(
		"POST",
		tokenURL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token exchange failed: %s", string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token not found in response")
	}

	return accessToken, nil
}

func (s *OAuthExchangeService) GetUserProfile(
	provider *models.AuthenticationProvider,
	accessToken string,
) (*UserProfile, error) {
	switch provider.Provider {
	case models.ProviderGoogle:
		return s.getGoogleUserProfile(accessToken)
	case models.ProviderGithub:
		return s.getGithubUserProfile(accessToken)
	case models.ProviderMicrosoft:
		return s.getMicrosoftUserProfile(accessToken)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider.Provider)
	}
}

func (s *OAuthExchangeService) getGoogleUserProfile(
	accessToken string,
) (*UserProfile, error) {
	req, err := http.NewRequest(
		"GET",
		"https://www.googleapis.com/oauth2/v2/userinfo",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user profile: %s", string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &UserProfile{
		ProviderUserID: result["id"].(string),
		Email:          result["email"].(string),
		Name:           result["name"].(string),
		AvatarURL:      result["picture"].(string),
	}, nil
}

func (s *OAuthExchangeService) getGithubUserProfile(
	accessToken string,
) (*UserProfile, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user profile: %s", string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	email := ""
	if result["email"] != nil {
		email = result["email"].(string)
	}

	if email == "" {
		email, err = s.getGithubPrimaryEmail(accessToken)
		if err != nil {
			return nil, err
		}
	}

	name := ""
	if result["name"] != nil {
		name = result["name"].(string)
	}

	avatarURL := ""
	if result["avatar_url"] != nil {
		avatarURL = result["avatar_url"].(string)
	}

	return &UserProfile{
		ProviderUserID: fmt.Sprintf("%v", result["id"]),
		Email:          email,
		Name:           name,
		AvatarURL:      avatarURL,
	}, nil
}

func (s *OAuthExchangeService) getGithubPrimaryEmail(
	accessToken string,
) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get user emails: %s", string(body))
	}

	var emails []map[string]interface{}
	err = json.Unmarshal(body, &emails)
	if err != nil {
		return "", err
	}

	for _, email := range emails {
		if primary, ok := email["primary"].(bool); ok && primary {
			return email["email"].(string), nil
		}
	}

	if len(emails) > 0 {
		return emails[0]["email"].(string), nil
	}

	return "", fmt.Errorf("no email found")
}

func (s *OAuthExchangeService) getMicrosoftUserProfile(
	accessToken string,
) (*UserProfile, error) {
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user profile: %s", string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	name := ""
	if result["displayName"] != nil {
		name = result["displayName"].(string)
	}

	return &UserProfile{
		ProviderUserID: result["id"].(string),
		Email:          result["mail"].(string),
		Name:           name,
		AvatarURL:      "",
	}, nil
}
