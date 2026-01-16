package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hisamafahri/securelogin/internal/usecase"
	"github.com/hisamafahri/securelogin/pkg/response"
)

type CallbackController struct {
	oauthCallbackUsecase *usecase.OAuthCallbackUsecase
}

func NewCallbackController(
	oauthCallbackUsecase *usecase.OAuthCallbackUsecase,
) *CallbackController {
	return &CallbackController{
		oauthCallbackUsecase: oauthCallbackUsecase,
	}
}

func (c *CallbackController) GoogleCallback(ctx *gin.Context) {
	c.handleCallback(ctx)
}

func (c *CallbackController) GithubCallback(ctx *gin.Context) {
	c.handleCallback(ctx)
}

func (c *CallbackController) MicrosoftCallback(ctx *gin.Context) {
	c.handleCallback(ctx)
}

func (c *CallbackController) handleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if code == "" {
		response.JSON(
			ctx,
			http.StatusBadRequest,
			"invalid_request",
			"missing code parameter",
		)
		return
	}

	if state == "" {
		response.JSON(
			ctx,
			http.StatusBadRequest,
			"invalid_request",
			"missing state parameter",
		)
		return
	}

	redirectURL, err := c.oauthCallbackUsecase.Execute(code, state)
	if err != nil {
		response.JSON(
			ctx,
			http.StatusInternalServerError,
			"server_error",
			err.Error(),
		)
		return
	}

	ctx.Redirect(http.StatusFound, redirectURL)
}
