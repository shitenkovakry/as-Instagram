package likes

import (
	"instagram/logger"

	"github.com/pkg/errors"
)

type DB interface {
	InsertForLike(idPhoto int, idUser int) error
	CountLikes(idPhoto int) (int, error)
	DeleteLike(photoID int, userID int) error
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

func (like *LikesManager) Count(idPhoto int) (int, error) {
	counted, err := like.db.CountLikes(idPhoto)
	if err != nil {
		return 0, errors.Wrap(err, "can not count likes")
	}

	return counted, nil
}

func (like *LikesManager) Delete(photoID int, userID int) error {
	err := like.db.DeleteLike(photoID, userID)
	if err != nil {
		return errors.Wrap(err, "can not delete like")
	}

	return nil
}
