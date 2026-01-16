package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hisamafahri/securelogin/internal/service"
)

type AccessTokenClaims struct {
	UserID   string
	ClientID string
	Scope    string
}

func RequireAccessToken(jwtService *service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":             "invalid_request",
				"error_description": "missing authorization header",
			})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":             "invalid_request",
				"error_description": "invalid authorization header format",
			})
			ctx.Abort()
			return
		}

		accessToken := parts[1]

		claims, err := jwtService.ValidateAccessToken(accessToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":             "invalid_token",
				"error_description": err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("access_token_claims", claims)
		ctx.Next()
	}
}
