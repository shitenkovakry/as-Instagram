package deletephoto

import (
	"encoding/json"
	"errors"
	"instagram/errordefs"
	"instagram/logger"
	models "instagram/models/photos"
	"io"
	"net/http"
)

type DeletePhoto struct {
	IdPhoto int `json:"photo_id"`
}

type PhotoActionsForHandlerDeletePhoto interface {
	Delete(photoID int) (*models.Photo, error)
}

type HandlerForDeletePhoto struct {
	log          logger.Logger
	photoActions PhotoActionsForHandlerDeletePhoto
}

func (handler *HandlerForDeletePhoto) prepareRequest(request *http.Request) (*models.Photo, error) {
	defer func() {
		if err := request.Body.Close(); err != nil {
			handler.log.Printf("cannot close body: %v", err)
		}
	}()

	body, err := io.ReadAll(request.Body)
	if err != nil {
		handler.log.Printf("cannot read body: %v", err)

		return nil, err
	}

	var deletePhotoFromClient *DeletePhoto

	if err := json.Unmarshal(body, &deletePhotoFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	deletedPhoto := &models.Photo{
		IDPhoto: deletePhotoFromClient.IdPhoto,
	}

	return deletedPhoto, nil
}

func (handler *HandlerForDeletePhoto) sendResponse(writer http.ResponseWriter, deletedPhoto *models.Photo) {
	deletedPhotoMarshaled, err := json.Marshal(deletedPhoto)
	if err != nil {
		handler.log.Printf("can not marshal deleted photo: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := writer.Write(deletedPhotoMarshaled); err != nil {
		handler.log.Printf("can not send to client deleted photo: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForDeletePhoto) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	shouldDeletePhoto, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("can not prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	deletedPhoto, err := handler.photoActions.Delete(shouldDeletePhoto.IDPhoto)
	if errors.Is(err, errordefs.ErrNoDocuments) {
		writer.WriteHeader(http.StatusNotFound)

		return
	}

	if err != nil {
		handler.log.Printf("can not delete photo: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(writer, deletedPhoto)
}

func NewHandlerForDeletePhoto(log logger.Logger, photoActions PhotoActionsForHandlerDeletePhoto) *HandlerForDeletePhoto {
	result := &HandlerForDeletePhoto{
		log:          log,
		photoActions: photoActions,
	}

	return result
}
