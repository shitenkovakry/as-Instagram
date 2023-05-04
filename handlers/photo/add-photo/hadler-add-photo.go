package addphoto

import (
	"encoding/json"
	"instagram/logger"
	models "instagram/models/photos"
	"io"
	"net/http"
)

type AddPhoto struct {
	IDUser int `json:"user_id"`
}

type PhotoActionsForHandlerAddPhoto interface {
	Add(photoOfUser *models.Photo) (*models.Photo, error)
}

type HandlerForAddPhoto struct {
	log          logger.Logger
	photoActions PhotoActionsForHandlerAddPhoto
}

func (handler *HandlerForAddPhoto) prepareRequest(request *http.Request) (*models.Photo, error) {
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

	var newPhotoFromClient *AddPhoto

	if err := json.Unmarshal(body, &newPhotoFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	newPhoto := &models.Photo{
		IDUser: newPhotoFromClient.IDUser,
	}

	return newPhoto, nil
}

func (handler *HandlerForAddPhoto) sendResponse(write http.ResponseWriter, addedPhoto *models.Photo) {
	addedPhotoMarshaled, err := json.Marshal(addedPhoto)
	if err != nil {
		handler.log.Printf("cannot marshal added photo: %v", err)
		write.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := write.Write(addedPhotoMarshaled); err != nil {
		handler.log.Printf("cannot send to client added photo: %v", err)
		write.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForAddPhoto) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	newPhoto, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("cannot prepare request: %v", err)
		write.WriteHeader(http.StatusBadRequest)

		return
	}

	createdPhoto, err := handler.photoActions.Add(newPhoto)
	if err != nil {
		handler.log.Printf("cannot create photo: %v", err)
		write.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(write, createdPhoto)
}

func NewHandlerForAddPhoto(log logger.Logger, photoActions PhotoActionsForHandlerAddPhoto) *HandlerForAddPhoto {
	result := &HandlerForAddPhoto{
		log:          log,
		photoActions: photoActions,
	}

	return result
}
