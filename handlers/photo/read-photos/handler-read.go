package readphotos

import (
	"encoding/json"
	"instagram/logger"
	models "instagram/models/photos"
	"io"
	"net/http"
)

type ReadPhotos struct {
	IDUser int `json:"user_id"`
}

type PhotosActionsForHandlerReadPhotos interface {
	ReadPhotos(userID int) (models.Photos, error)
}

type HandlerForReadPhotos struct {
	log           logger.Logger
	photosActions PhotosActionsForHandlerReadPhotos
}

func NewHandlerForReadPhotos(log logger.Logger, photosActions PhotosActionsForHandlerReadPhotos) *HandlerForReadPhotos {
	result := &HandlerForReadPhotos{
		log:           log,
		photosActions: photosActions,
	}

	return result
}

func (handler *HandlerForReadPhotos) prepareRequest(request *http.Request) (*models.Photo, error) {
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

	var readPhotoFromClient *ReadPhotos

	if err := json.Unmarshal(body, &readPhotoFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	readPhoto := &models.Photo{
		IDUser: readPhotoFromClient.IDUser,
	}

	return readPhoto, nil
}

func (handler *HandlerForReadPhotos) sendResponse(writer http.ResponseWriter, readPhotos models.Photos) {
	data, err := json.Marshal(readPhotos)
	if err != nil {
		handler.log.Printf("can not marshal photos: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := writer.Write(data); err != nil {
		handler.log.Printf("can not writer photos: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForReadPhotos) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	shouldReadPhotos, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("can not prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	photos, err := handler.photosActions.ReadPhotos(shouldReadPhotos.IDUser)
	if err != nil {
		handler.log.Printf("can not read photos: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(writer, photos)
}
