package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
	"github.com/hisamafahri/securelogin/internal/service"
)

type WellKnownController struct {
	jwtService *service.JWTService
}

func NewWellKnownController(
	jwtService *service.JWTService,
) *WellKnownController {
	return &WellKnownController{
		jwtService: jwtService,
	}
}

func (c *WellKnownController) OpenIDConfiguration(ctx *gin.Context) {
	baseURL := os.Getenv("SYSTEM_BASE_URL")

	config := map[string]interface{}{
		"issuer":                   baseURL,
		"authorization_endpoint":   baseURL + "/authorize",
		"token_endpoint":           baseURL + "/oauth/token",
		"userinfo_endpoint":        baseURL + "/oauth/userinfo",
		"jwks_uri":                 baseURL + "/.well-known/jwks.json",
		"revocation_endpoint":      baseURL + "/oauth/revoke",
		"response_types_supported": []string{"code"},
		"grant_types_supported": []string{
			"authorization_code",
			"refresh_token",
		},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
		"token_endpoint_auth_methods_supported": []string{"client_secret_post"},
		"scopes_supported": []string{
			"openid",
			"profile",
			"email",
		},
		"claims_supported": []string{
			"sub",
			"iss",
			"aud",
			"exp",
			"iat",
			"name",
			"email",
			"email_verified",
			"picture",
		},
		"code_challenge_methods_supported": []string{"S256"},
	}

	ctx.JSON(http.StatusOK, config)
}

func (c *WellKnownController) JWKS(ctx *gin.Context) {
	publicKeys := c.jwtService.GetPublicKeys()

	jwks := jose.JSONWebKeySet{
		Keys: publicKeys,
	}

	ctx.JSON(http.StatusOK, jwks)
}
