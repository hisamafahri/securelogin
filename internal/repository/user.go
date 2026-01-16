package repository

import (
	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/infrastructure/pgsql/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindOrCreate(
	user *models.User,
) (*models.User, error) {
	err := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "application_id"},
			{Name: "provider_user_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"email",
			"name",
			"avatar_url",
			"updated_at",
		}),
	}).Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
