package repository

import (
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/models"
)

type UserRepository interface {
	Create(*models.User) error
	FindByEmail(email string) (*models.User, error)
	Delete(*models.User) error
}

func NewUserRepository(dbWrapper *models.DBWrapper) UserRepository {
	return &userRepository{DBWrapper: dbWrapper}
}

type userRepository struct {
	DBWrapper *models.DBWrapper
}

func (u *userRepository) Create(newUser *models.User) error {
	return u.DBWrapper.PG_DBConnection.Create(newUser).Error
}

func (u *userRepository) FindByEmail(email string) (*models.User, error) {
	var result models.User
	err := u.DBWrapper.PG_DBConnection.Where("email = ?", email).First(&result).Error
	return &result, err
}

func (u *userRepository) Delete(removeUser *models.User) error {
	err := u.DBWrapper.PG_DBConnection.Delete(removeUser).Error
	return err
}
