package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"github.com/hisamafahri/securelogin/pkg/utils"
	"gorm.io/gorm"
)

var ErrInvalidClientCredentials = errors.New("invalid client credentials")

type ApplicationService struct {
	db *gorm.DB
}

func NewApplicationService(db *gorm.DB) *ApplicationService {
	return &ApplicationService{db: db}
}

func (s *ApplicationService) GetByID(
	id uuid.UUID,
) (*models.Application, error) {
	var app models.Application
	err := s.db.Where("id = ?", id).First(&app).Error
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (s *ApplicationService) GetByClientID(
	clientID string,
) (*models.Application, error) {
	var app models.Application
	err := s.db.Where("client_id = ?", clientID).First(&app).Error
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (s *ApplicationService) IsRedirectURIValid(
	app *models.Application,
	redirectURI string,
) bool {
	parsedURI, err := utils.SafelyParseURL(redirectURI)
	if err != nil {
		return false
	}

	if parsedURI.Origin == nil {
		return false
	}

	normalizedInput := *parsedURI.Origin + parsedURI.Pathname

	for _, registeredURI := range app.RedirectURIs {
		trimmedURI := strings.Trim(registeredURI, "'\"")

		parsedRegistered, err := utils.SafelyParseURL(trimmedURI)
		if err != nil {
			continue
		}

		if parsedRegistered.Origin == nil {
			continue
		}

		normalizedRegistered := *parsedRegistered.Origin + parsedRegistered.Pathname

		if normalizedInput == normalizedRegistered {
			return true
		}
	}

	return false
}

func (s *ApplicationService) ValidateClientCredentials(
	clientID string,
	clientSecret string,
) (*models.Application, error) {
	app, err := s.GetByClientID(clientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidClientCredentials
		}
		return nil, err
	}

	if app.ClientSecret != clientSecret {
		return nil, ErrInvalidClientCredentials
	}

	return app, nil
}
