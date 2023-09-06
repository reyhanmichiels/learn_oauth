package repository

import (
	"learn_oauth/domain"

	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *domain.User) error
	UpdateUserReturnAffectedRow(user *domain.User) (int64, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) CreateUser(user *domain.User) error {
	err := repo.db.Create(&user).Error
	return err
}

func (repo *UserRepository) UpdateUserReturnAffectedRow(user *domain.User) (int64, error) {
	result := repo.db.Model(user).Where("email = ?", user.Email).Updates(user)

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
