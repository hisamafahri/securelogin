package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/internal/service"
	"github.com/hisamafahri/securelogin/internal/usecase"
)

type UserinfoController struct {
	userinfoUsecase *usecase.UserinfoUsecase
}

func NewUserinfoController(
	userinfoUsecase *usecase.UserinfoUsecase,
) *UserinfoController {
	return &UserinfoController{
		userinfoUsecase: userinfoUsecase,
	}
}

func (c *UserinfoController) Userinfo(ctx *gin.Context) {
	claimsInterface, exists := ctx.Get("access_token_claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":             "invalid_token",
			"error_description": "missing token claims",
		})
		return
	}

	claims, ok := claimsInterface.(*service.AccessTokenClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":             "server_error",
			"error_description": "failed to parse token claims",
		})
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":             "server_error",
			"error_description": "invalid user id in token",
		})
		return
	}

	user, err := c.userinfoUsecase.GetUserInfo(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":             "not_found",
			"error_description": "user not found",
		})
		return
	}

	response := gin.H{
		"sub":            user.ID.String(),
		"email":          user.Email,
		"email_verified": true,
	}

	if user.Name != nil {
		response["name"] = *user.Name
	}

	if user.AvatarURL != nil {
		response["picture"] = *user.AvatarURL
	}

	ctx.JSON(http.StatusOK, response)
}
