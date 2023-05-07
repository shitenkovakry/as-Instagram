package photos

import (
	"instagram/logger"
	models "instagram/models/photos"

	"github.com/pkg/errors"
)

type DB interface {
	ReadPhoto(userID int) (models.Photos, error)
	InsertForPhoto(photoOfUser *models.Photo) (*models.Photo, error)
	DeletePhoto(photoID int) (*models.Photo, error)
}

type PhotosManager struct {
	db  DB
	log logger.Logger
}

func New(log logger.Logger, db DB) *PhotosManager {
	return &PhotosManager{
		log: log,
		db:  db,
	}
}

func (photos *PhotosManager) Read(userID int) (models.Photos, error) {
	read, err := photos.db.ReadPhoto(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "can not read")
	}

	return read, nil
}

func (photo *PhotosManager) Add(photoOfUser *models.Photo) (*models.Photo, error) {
	insertedPhoto, err := photo.db.InsertForPhoto(photoOfUser)
	if err != nil {
		return nil, errors.Wrap(err, "cannot add photo")
	}

	return insertedPhoto, nil
}

func (photo *PhotosManager) Delete(photoID int) (*models.Photo, error) {
	deletedPhoto, err := photo.db.DeletePhoto(photoID)
	if err != nil {
		return nil, errors.Wrap(err, "can not delete photo")
	}

	return deletedPhoto, nil
}
