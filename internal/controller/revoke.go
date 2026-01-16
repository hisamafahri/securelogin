package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hisamafahri/securelogin/internal/repository"
	"github.com/hisamafahri/securelogin/internal/service"
	"github.com/hisamafahri/securelogin/internal/usecase"
	"github.com/hisamafahri/securelogin/pkg/response"
)

type RevokeRequest struct {
	ClientID     string `form:"client_id"     validate:"required"`
	ClientSecret string `form:"client_secret" validate:"required"`
	RefreshToken string `form:"refresh_token" validate:"required"`
}

type RevokeController struct {
	tokenRevokeUsecase *usecase.TokenRevokeUsecase
	validate           *validator.Validate
}

func NewRevokeController(
	tokenRevokeUsecase *usecase.TokenRevokeUsecase,
	validate *validator.Validate,
) *RevokeController {
	return &RevokeController{
		tokenRevokeUsecase: tokenRevokeUsecase,
		validate:           validate,
	}
}

func (ctrl *RevokeController) Revoke(c *gin.Context) {
	var req RevokeRequest

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

	err := ctrl.tokenRevokeUsecase.RevokeToken(
		req.ClientID,
		req.ClientSecret,
		req.RefreshToken,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidClientCredentials) {
			response.JSON(
				c,
				http.StatusUnauthorized,
				"invalid_client",
				"invalid client credentials",
			)
			return
		}
		if errors.Is(err, repository.ErrTokenNotFound) {
			response.JSON(
				c,
				http.StatusBadRequest,
				"invalid_grant",
				"token not found or already revoked",
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

	c.Status(http.StatusOK)
}
