package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/internal/service"
	"github.com/hisamafahri/securelogin/internal/usecase"
	"github.com/hisamafahri/securelogin/internal/view"
	"github.com/hisamafahri/securelogin/pkg/response"
)

type SigninController struct {
	appService              *service.ApplicationService
	providerService         *service.AuthenticationProviderService
	authReqService          *service.AuthenticationRequestService
	signinIdentifierUsecase *usecase.SigninIdentifierUsecase
}

func NewSigninController(
	appService *service.ApplicationService,
	providerService *service.AuthenticationProviderService,
	authReqService *service.AuthenticationRequestService,
	signinIdentifierUsecase *usecase.SigninIdentifierUsecase,
) *SigninController {
	return &SigninController{
		appService:              appService,
		providerService:         providerService,
		authReqService:          authReqService,
		signinIdentifierUsecase: signinIdentifierUsecase,
	}
}

func (ctrl *SigninController) Signin(c *gin.Context) {
	rayID := c.GetString("ray_id")
	requestID := c.Query("request_id")
	if requestID == "" {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			"request_id is required",
		)
		return
	}

	requestUUID, err := uuid.Parse(requestID)
	if err != nil {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			"request_id is invalid",
		)
		return
	}

	authReq, err := ctrl.authReqService.GetByID(requestUUID)
	if err != nil {
		response.JSON(
			c,
			http.StatusNotFound,
			"request_not_found",
			"authentication request not found",
		)
		return
	}

	if authReq.CompletedAt != nil {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			"authentication request has already been completed",
		)
		return
	}

	if authReq.ExpiresAt.Before(time.Now()) {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			"authentication request has expired",
		)
		return
	}

	ipAddress := c.ClientIP()

	app, err := ctrl.appService.GetByID(authReq.ApplicationID)
	if err != nil {
		response.JSON(
			c,
			http.StatusNotFound,
			"application_not_found",
			"application not found",
		)
		return
	}

	providers, err := ctrl.providerService.GetProviders(
		app.ID,
	)
	if err != nil {
		response.JSON(
			c,
			http.StatusInternalServerError,
			"internal_server_error",
			"failed to fetch providers",
		)
		return
	}

	signinView := view.NewSigninView(
		requestID,
		rayID,
		ipAddress,
		app.Name,
		providers,
	)
	signinView.Render(c)
}

func (ctrl *SigninController) SigninIdentifier(c *gin.Context) {
	providerIDStr := c.PostForm("provider_id")
	requestIDStr := c.PostForm("request_id")

	if providerIDStr == "" {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			"provider_id is required",
		)
		return
	}

	if requestIDStr == "" {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			"request_id is required",
		)
		return
	}

	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			"provider_id is invalid",
		)
		return
	}

	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			"request_id is invalid",
		)
		return
	}

	authReq, err := ctrl.authReqService.GetByID(requestID)
	if err != nil {
		response.JSON(
			c,
			http.StatusNotFound,
			"request_not_found",
			"authentication request not found",
		)
		return
	}

	if authReq.CompletedAt != nil {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			"authentication request has already been completed",
		)
		return
	}

	if authReq.ExpiresAt.Before(time.Now()) {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			"authentication request has expired",
		)
		return
	}

	oauthURL, err := ctrl.signinIdentifierUsecase.Execute(providerID, requestID)
	if err != nil {
		response.JSON(
			c,
			http.StatusBadRequest,
			"invalid_request",
			err.Error(),
		)
		return
	}

	c.Redirect(http.StatusFound, oauthURL)
}
