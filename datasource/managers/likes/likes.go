package likes

import (
	"instagram/logger"

	"github.com/pkg/errors"
)

type DB interface {
	InsertForLike(idPhoto int, idUser int) error
}

type LikesManager struct {
	log logger.Logger
	db  DB
}

func New(log logger.Logger, db DB) *LikesManager {
	return &LikesManager{
		log: log,
		db:  db,
	}
}

func (like *LikesManager) Add(idPhoto int, idUser int) error {
	err := like.db.InsertForLike(idPhoto, idUser)
	if err != nil {
		return errors.Wrap(err, "cannot add like")
	}

	return nil
}
