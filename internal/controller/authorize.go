package controller

import (
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hisamafahri/securelogin/internal/usecase"
	"github.com/hisamafahri/securelogin/pkg/response"
)

type AuthorizeRequest struct {
	ResponseType        string  `form:"response_type"         validate:"required,oneof=code token"`
	ClientID            string  `form:"client_id"             validate:"required"`
	RedirectURI         string  `form:"redirect_uri"          validate:"required,url"`
	State               *string `form:"state"`
	Scope               *string `form:"scope"`
	CodeChallenge       *string `form:"code_challenge"`
	CodeChallengeMethod *string `form:"code_challenge_method" validate:"omitempty,oneof=S256"`
}

type AuthorizeController struct {
	authorizeUsecase *usecase.AuthorizeUsecase
	validate         *validator.Validate
}

func NewAuthorizeController(
	authorizeUsecase *usecase.AuthorizeUsecase,
	validate *validator.Validate,
) *AuthorizeController {
	return &AuthorizeController{
		authorizeUsecase: authorizeUsecase,
		validate:         validate,
	}
}

func (ctrl *AuthorizeController) Authorize(c *gin.Context) {
	var req AuthorizeRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	if err := ctrl.validate.Struct(&req); err != nil {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_payload",
			response.FormatValidationError(err),
		)
		return
	}

	if req.Scope != nil {
		normalizedScope := strings.Join(strings.Fields(*req.Scope), " ")
		validScopes := []string{
			"offline_access",
			"openid",
			"offline_access openid",
			"openid offline_access",
		}
		isValid := slices.Contains(validScopes, normalizedScope)
		if !isValid {
			response.JSON(
				c,
				http.StatusBadRequest,
				"invalid_scope",
				"scope must be one of: offline_access, openid, or both",
			)
			return
		}
	}

	app, err := ctrl.authorizeUsecase.ValidateAndGetApplication(
		req.ClientID,
		req.RedirectURI,
	)
	if err != nil {
		if errors.Is(err, usecase.ErrClientNotFound) {
			response.JSON(
				c,
				http.StatusNotFound,
				"application_not_found",
				"application is not found",
			)
			return
		}
		if errors.Is(err, usecase.ErrInvalidRedirectURI) {
			response.JSON(
				c,
				http.StatusBadRequest,
				"invalid_redirect_uri",
				"redirect_uri not registered for this application",
			)
			return
		}
		response.JSON(
			c,
			http.StatusInternalServerError,
			"internal_server_error",
			"internal server error",
		)
		return
	}

	// NOTE:
	// split scope into array, trim it, and filter unique scopes
	scopes := []string{}
	if req.Scope != nil {
		scopeMap := make(map[string]struct{})
		for scope := range strings.SplitSeq(*req.Scope, " ") {
			trimmed := strings.TrimSpace(scope)
			if trimmed != "" {
				scopeMap[trimmed] = struct{}{}
			}
		}
		for scope := range scopeMap {
			scopes = append(scopes, scope)
		}
	}

	authReq, err := ctrl.authorizeUsecase.CreateAuthenticationRequest(
		app.ID,
		req.ResponseType,
		req.RedirectURI,
		req.State,
		req.CodeChallenge,
		req.CodeChallengeMethod,
		scopes,
	)
	if err != nil {
		response.JSON(
			c,
			http.StatusInternalServerError,
			"internal_server_error",
			"internal server error",
		)
		return
	}

	params := url.Values{}
	params.Add("request_id", authReq.ID.String())

	c.Redirect(http.StatusFound, "/signin?"+params.Encode())
}
