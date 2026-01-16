package usecase

import (
	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"github.com/hisamafahri/securelogin/internal/service"
)

type UserinfoUsecase struct {
	userService *service.UserService
}

func NewUserinfoUsecase(userService *service.UserService) *UserinfoUsecase {
	return &UserinfoUsecase{
		userService: userService,
	}
}

func (u *UserinfoUsecase) GetUserInfo(userID uuid.UUID) (*models.User, error) {
	return u.userService.GetByID(userID)
}
