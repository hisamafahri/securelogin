package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hisamafahri/securelogin/internal/middleware"
	"github.com/hisamafahri/securelogin/internal/service"
)

func RegisterRoutes(
	r *gin.Engine,
	authorizeController *AuthorizeController,
	signinController *SigninController,
	callbackController *CallbackController,
	tokenController *TokenController,
	revokeController *RevokeController,
	wellknownController *WellKnownController,
	userinfoController *UserinfoController,
	jwtService *service.JWTService,
) {
	r.GET("/ping", Ping)
	r.GET("/authorize", authorizeController.Authorize)
	r.GET("/signin", signinController.Signin)
	r.POST("/signin/identifier", signinController.SigninIdentifier)
	r.GET("/callback/google", callbackController.GoogleCallback)
	r.GET("/callback/github", callbackController.GithubCallback)
	r.GET("/callback/microsoft", callbackController.MicrosoftCallback)
	r.POST("/oauth/token", tokenController.Token)
	r.POST("/oauth/revoke", revokeController.Revoke)

	r.GET(
		"/.well-known/openid-configuration",
		wellknownController.OpenIDConfiguration,
	)
	r.GET("/.well-known/jwks.json", wellknownController.JWKS)

	r.GET(
		"/oauth/userinfo",
		middleware.RequireAccessToken(jwtService),
		userinfoController.Userinfo,
	)
}
