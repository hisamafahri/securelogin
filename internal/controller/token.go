package controller

import (
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hisamafahri/securelogin/internal/service"
	"github.com/hisamafahri/securelogin/internal/usecase"
	"github.com/hisamafahri/securelogin/pkg/response"
)

type TokenRequest struct {
	GrantType    string `form:"grant_type"    validate:"required"`
	ClientID     string `form:"client_id"     validate:"required"`
	ClientSecret string `form:"client_secret" validate:"required"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	CodeVerifier string `form:"code_verifier"`
	RefreshToken string `form:"refresh_token"`
	Scope        string `form:"scope"`
}

type TokenResponse struct {
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type TokenController struct {
	tokenExchangeUsecase *usecase.TokenExchangeUsecase
	jwtService           *service.JWTService
	validate             *validator.Validate
}

func NewTokenController(
	tokenExchangeUsecase *usecase.TokenExchangeUsecase,
	jwtService *service.JWTService,
	validate *validator.Validate,
) *TokenController {
	return &TokenController{
		tokenExchangeUsecase: tokenExchangeUsecase,
		jwtService:           jwtService,
		validate:             validate,
	}
}

func (ctrl *TokenController) Token(c *gin.Context) {
	var req TokenRequest

	if err := c.ShouldBind(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	if err := ctrl.validate.Struct(&req); err != nil {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			response.FormatValidationError(err),
		)
		return
	}

	switch req.GrantType {
	case "authorization_code":
		ctrl.handleAuthorizationCode(c, req)
	case "refresh_token":
		ctrl.handleRefreshToken(c, req)
	default:
		response.JSON(
			c,
			http.StatusBadRequest,
			"unsupported_grant_type",
			"grant type must be authorization_code or refresh_token",
		)
	}
}

func (ctrl *TokenController) handleAuthorizationCode(
	c *gin.Context,
	req TokenRequest,
) {
	if req.Code == "" {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			"code is required for authorization_code grant type",
		)
		return
	}

	session, user, authReq, err := ctrl.tokenExchangeUsecase.ExchangeCodeForToken(
		req.GrantType,
		req.ClientID,
		req.ClientSecret,
		req.Code,
		req.RedirectURI,
		req.CodeVerifier,
	)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidGrantType) {
			response.JSON(
				c,
				http.StatusBadRequest,
				"unsupported_grant_type",
				"grant type must be authorization_code",
			)
			return
		}
		if errors.Is(err, service.ErrInvalidClientCredentials) {
			response.JSON(
				c,
				http.StatusUnauthorized,
				"invalid_client",
				"invalid client credentials",
			)
			return
		}
		if errors.Is(err, service.ErrAuthorizationCodeNotFound) {
			response.JSON(
				c,
				http.StatusBadRequest,
				"invalid_grant",
				"authorization code not found",
			)
			return
		}
		if errors.Is(err, service.ErrAuthorizationCodeExpired) {
			response.JSON(
				c,
				http.StatusBadRequest,
				"invalid_grant",
				"authorization code expired",
			)
			return
		}
		if errors.Is(err, usecase.ErrAuthorizationCodeUsed) {
			response.JSON(
				c,
				http.StatusBadRequest,
				"invalid_grant",
				"authorization code already used",
			)
			return
		}
		if errors.Is(err, usecase.ErrRedirectURIMismatch) {
			response.JSON(
				c,
				http.StatusBadRequest,
				"invalid_grant",
				"redirect_uri does not match",
			)
			return
		}
		if errors.Is(err, usecase.ErrCodeVerifierRequired) {
			response.JSON(
				c,
				http.StatusBadRequest,
				"invalid_request",
				"code_verifier is required for PKCE",
			)
			return
		}
		if errors.Is(err, usecase.ErrInvalidCodeVerifier) {
			response.JSON(
				c,
				http.StatusBadRequest,
				"invalid_grant",
				"invalid code_verifier",
			)
			return
		}
		response.JSON(
			c,
			http.StatusInternalServerError,
			"server_error",
			"internal server error",
		)
		return
	}

	expiresIn := int(session.ExpiresAt.Sub(session.CreatedAt).Seconds())
	tokenExpiry := time.Duration(expiresIn) * time.Second

	scope := strings.Join(authReq.Scopes, " ")

	accessToken, err := ctrl.jwtService.GenerateAccessToken(
		user.ID,
		req.ClientID,
		scope,
		tokenExpiry,
	)
	if err != nil {
		response.JSON(
			c,
			http.StatusInternalServerError,
			"server_error",
			"failed to generate access token",
		)
		return
	}

	var idToken string
	hasOpenID := slices.Contains(authReq.Scopes, "openid")
	hasOfflineAccess := slices.Contains(authReq.Scopes, "offline_access")

	if hasOpenID {
		idToken, err = ctrl.jwtService.GenerateIDToken(
			user.ID,
			req.ClientID,
			user.Email,
			user.Name,
			user.AvatarURL,
			tokenExpiry,
		)
		if err != nil {
			response.JSON(
				c,
				http.StatusInternalServerError,
				"server_error",
				"failed to generate id token",
			)
			return
		}
	}

	refreshToken := ""
	if hasOfflineAccess {
		refreshToken = session.Token
	}

	c.JSON(http.StatusOK, TokenResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		IDToken:      idToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	})
}

func (ctrl *TokenController) handleRefreshToken(
	c *gin.Context,
	req TokenRequest,
) {
	if req.RefreshToken == "" {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			"refresh_token is required for refresh_token grant type",
		)
		return
	}

	session, user, _, scopes, err := ctrl.tokenExchangeUsecase.RefreshAccessToken(
		req.GrantType,
		req.ClientID,
		req.ClientSecret,
		req.RefreshToken,
		req.Scope,
	)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidGrantType) {
			response.JSON(
				c,
				http.StatusBadRequest,
				"unsupported_grant_type",
				"grant type must be refresh_token",
			)
			return
		}
		if errors.Is(err, service.ErrInvalidClientCredentials) {
			response.JSON(
				c,
				http.StatusUnauthorized,
				"invalid_client",
				"invalid client credentials",
			)
			return
		}
		if errors.Is(err, service.ErrRefreshTokenNotFound) {
			response.JSON(
				c,
				http.StatusBadRequest,
				"invalid_grant",
				"refresh token not found",
			)
			return
		}
		if errors.Is(err, service.ErrRefreshTokenExpired) {
			response.JSON(
				c,
				http.StatusBadRequest,
				"invalid_grant",
				"refresh token expired",
			)
			return
		}
		if errors.Is(err, usecase.ErrInvalidScope) {
			response.JSON(
				c,
				http.StatusBadRequest,
				"invalid_scope",
				"requested scope is invalid or exceeds originally granted scope",
			)
			return
		}
		response.JSON(
			c,
			http.StatusInternalServerError,
			"server_error",
			"internal server error",
		)
		return
	}

	expiresIn := int(session.ExpiresAt.Sub(session.CreatedAt).Seconds())
	tokenExpiry := time.Duration(expiresIn) * time.Second

	scope := strings.Join(scopes, " ")

	accessToken, err := ctrl.jwtService.GenerateAccessToken(
		user.ID,
		req.ClientID,
		scope,
		tokenExpiry,
	)
	if err != nil {
		response.JSON(
			c,
			http.StatusInternalServerError,
			"server_error",
			"failed to generate access token",
		)
		return
	}

	var idToken string
	hasOpenID := slices.Contains(scopes, "openid")
	hasOfflineAccess := slices.Contains(scopes, "offline_access")

	if hasOpenID {
		idToken, err = ctrl.jwtService.GenerateIDToken(
			user.ID,
			req.ClientID,
			user.Email,
			user.Name,
			user.AvatarURL,
			tokenExpiry,
		)
		if err != nil {
			response.JSON(
				c,
				http.StatusInternalServerError,
				"server_error",
				"failed to generate id token",
			)
			return
		}
	}

	refreshToken := ""
	if hasOfflineAccess {
		refreshToken = session.Token
	}

	c.JSON(http.StatusOK, TokenResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		IDToken:      idToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	})
}
