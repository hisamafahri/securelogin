package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"github.com/hisamafahri/securelogin/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) FindOrCreateUser(
	applicationID uuid.UUID,
	providerID uuid.UUID,
	providerUserID string,
	email string,
	name *string,
	avatarURL *string,
) (*models.User, error) {
	user := &models.User{
		ApplicationID:  applicationID,
		ProviderID:     providerID,
		ProviderUserID: providerUserID,
		Email:          email,
		Name:           name,
		AvatarURL:      avatarURL,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return s.repo.FindOrCreate(user)
}

func (s *UserService) GetByID(id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(id)
}
