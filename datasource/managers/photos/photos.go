package photos

import (
	"fmt"
	"instagram/logger"
	models "instagram/models/photos"
	"os"

	"github.com/pkg/errors"
)

type DB interface {
	ReadPhotos(userID int) (models.Photos, error)
	ReadPhoto(idUser int, idPhoto int) (*models.Photo, error)
	InsertForPhoto(userID int, path string) (*models.Photo, error)
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

func (photos *PhotosManager) ReadPhotos(userID int) (models.Photos, error) {
	read, err := photos.db.ReadPhotos(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "can not read")
	}

	return read, nil
}

const (
	basePath = "/Users/kryshitenkova/Public/photos"
)

func (photos *PhotosManager) ReadPhoto(idUser int, idPhoto int) ([]byte, error) {
	read, err := photos.db.ReadPhoto(idUser, idPhoto)
	if err != nil {
		return nil, errors.Wrapf(err, "can not read from DB")
	}

	path := fmt.Sprintf("%s/%s", basePath, read.Path)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "can not read the photo from the file system")
	}

	return data, nil
}

func (photo *PhotosManager) Add(userID int, photoContent []byte, photoFilename string) (*models.Photo, error) {
	path := fmt.Sprintf("%s/%s", basePath, photoFilename)

	if err := os.WriteFile(path, photoContent, os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "cannot save photo content")
	}

	insertedPhoto, err := photo.db.InsertForPhoto(userID, photoFilename)
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
