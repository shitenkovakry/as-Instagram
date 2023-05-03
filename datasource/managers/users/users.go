package users

import (
	"instagram/logger"
	models "instagram/models/users"

	"github.com/pkg/errors"
)

type DB interface {
	Insert(user *models.UserRegistration) (*models.UserRegistration, error)
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

func (registration *UsersManager) Create(newUser *models.UserRegistration) (*models.UserRegistration, error) {
	insertedUser, err := registration.db.Insert(newUser)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create user")
	}

	return insertedUser, nil
}
