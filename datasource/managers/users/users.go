package users

import (
	"instagram/logger"
	models "instagram/models/users"

	"github.com/pkg/errors"
)

type DB interface {
	Insert(user *models.UserRegistration) (*models.UserRegistration, error)
	UpdateName(userID int, newName string) (*models.UserRegistration, error)
	UpdateEmail(userID int, newEmail string) (*models.UserRegistration, error)
}

type UsersManager struct {
	db  DB
	log logger.Logger
}

func New(log logger.Logger, db DB) *UsersManager {
	return &UsersManager{
		log: log,
		db:  db,
	}
}

func (user *UsersManager) Create(newUser *models.UserRegistration) (*models.UserRegistration, error) {
	insertedUser, err := user.db.Insert(newUser)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create user")
	}

	return insertedUser, nil
}

func (user *UsersManager) UpdateByName(userID int, newName string) (*models.UserRegistration, error) {
	updatedName, err := user.db.UpdateName(userID, newName)
	if err != nil {
		return nil, errors.Wrap(err, "can not update users name")
	}

	return updatedName, nil
}

func (user *UsersManager) UpdateByEmail(userID int, newEmail string) (*models.UserRegistration, error) {
	updatedEmail, err := user.db.UpdateEmail(userID, newEmail)
	if err != nil {
		return nil, errors.Wrap(err, "can not update users email")
	}

	return updatedEmail, nil
}
