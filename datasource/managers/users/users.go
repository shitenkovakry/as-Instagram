package users

import (
	"instagram/logger"
	models "instagram/models/users"

	"github.com/pkg/errors"
)

type DB interface {
	Insert(user *models.UserRegistration) (*models.UserRegistration, error)
	Update(userID int, newName string) (*models.UserRegistration, error)
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

func (user *UsersManager) Update(userID int, newName string) (*models.UserRegistration, error) {
	updatedName, err := user.db.Update(userID, newName)
	if err != nil {
		return nil, errors.Wrap(err, "can not update users name")
	}

	return updatedName, nil
}
