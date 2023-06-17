package photos

import (
	"bytes"
	"fmt"
	"instagram/logger"
	models "instagram/models/photos"
	"io"
	"log"
	"net/http"

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

func (photos *PhotosManager) ReadPhoto(idUser int, idPhoto int) ([]byte, error) {
	photo, err := photos.db.ReadPhoto(idUser, idPhoto)
	if err != nil {
		return nil, errors.Wrapf(err, "can not read from DB")
	}

	data, err := readFileFromServer3(photo.Path, idUser)
	if err != nil {
		return nil, errors.Wrapf(err, "can not read the photo from the external server")
	}

	return data, nil
}

func (photo *PhotosManager) Add(userID int, photoContent []byte, photoFilename string) (*models.Photo, error) {
	if err := sendFileToServer3(photoFilename, photoContent, userID); err != nil {
		return nil, err
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

func sendFileToServer3(photoFilename string, photoContent []byte, userID int) error {
	response, err := http.Post(
		fmt.Sprintf("http://server3:8080/upload/%d/%s", userID, photoFilename),
		"application/octet-stream",
		bytes.NewBuffer(photoContent),
	)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("something went wrong")
	}

	return nil
}

func readFileFromServer3(photoFilename string, userID int) ([]byte, error) {
	response, err := http.Get(
		fmt.Sprintf("http://server3:8080/download/%d/%s", userID, photoFilename),
	)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("something went wrong")
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("cannot close body in readFileFromServer3:%v", err)
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
