package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql"
	"github.com/hisamafahri/securelogin/internal/controller"
	"github.com/hisamafahri/securelogin/internal/middleware"
	"github.com/hisamafahri/securelogin/internal/repository"
	"github.com/hisamafahri/securelogin/internal/service"
	"github.com/hisamafahri/securelogin/internal/usecase"
)

func NewServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RayID())
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}

	// validator
	validate := validator.New()

	// repositories
	authProviderRepo := repository.NewAuthenticationProviderRepository(pgsql.DB)
	authReqRepo := repository.NewAuthenticationRequestRepository(pgsql.DB)
	userRepo := repository.NewUserRepository(pgsql.DB)
	authCodeRepo := repository.NewAuthorizationCodeRepository(pgsql.DB)
	sessionRepo := repository.NewSessionRepository(pgsql.DB)

	// services
	appService := service.NewApplicationService(pgsql.DB)
	authProviderService := service.NewAuthenticationProviderService(
		authProviderRepo,
	)
	authReqService := service.NewAuthenticationRequestService(authReqRepo)
	userService := service.NewUserService(userRepo)
	authCodeService := service.NewAuthorizationCodeService(authCodeRepo)
	sessionService := service.NewSessionService(sessionRepo)
	oauthExchangeService := service.NewOAuthExchangeService()

	issuer := os.Getenv("SYSTEM_BASE_URL")
	if issuer == "" {
		log.Fatal("SYSTEM_BASE_URL environment variable is required")
	}

	keyStorePath := os.Getenv("JWT_KEY_STORE_PATH")
	if keyStorePath == "" {
		keyStorePath = "./data/jwt-keys"
	}

	jwtService, err := service.NewJWTService(issuer, keyStorePath)
	if err != nil {
		log.Fatalf("failed to initialize JWT service: %v", err)
	}

	// usecases
	authorizeUsecase := usecase.NewAuthorizeUsecase(appService, authReqService)
	signinIdentifierUsecase := usecase.NewSigninIdentifierUsecase(
		authReqService,
		authProviderService,
	)
	oauthCallbackUsecase := usecase.NewOAuthCallbackUsecase(
		authReqService,
		authProviderService,
		oauthExchangeService,
		userService,
		authCodeService,
		sessionService,
	)
	tokenExchangeUsecase := usecase.NewTokenExchangeUsecase(
		appService,
		authCodeService,
		sessionService,
	)
	tokenRevokeUsecase := usecase.NewTokenRevokeUsecase(
		appService,
		sessionService,
	)
	userinfoUsecase := usecase.NewUserinfoUsecase(userService)

	// controllers
	authorizeController := controller.NewAuthorizeController(
		authorizeUsecase,
		validate,
	)
	signinController := controller.NewSigninController(
		appService,
		authProviderService,
		authReqService,
		signinIdentifierUsecase,
	)
	callbackController := controller.NewCallbackController(oauthCallbackUsecase)
	tokenController := controller.NewTokenController(
		tokenExchangeUsecase,
		jwtService,
		validate,
	)
	revokeController := controller.NewRevokeController(
		tokenRevokeUsecase,
		validate,
	)
	wellknownController := controller.NewWellKnownController(jwtService)
	userinfoController := controller.NewUserinfoController(userinfoUsecase)

	controller.RegisterRoutes(
		r,
		authorizeController,
		signinController,
		callbackController,
		tokenController,
		revokeController,
		wellknownController,
		userinfoController,
		jwtService,
	)
	return r
}
