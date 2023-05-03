package readphoto

import (
	"encoding/json"
	"instagram/logger"
	models "instagram/models/photos"
	"net/http"
)

type PhotoActionsForHandlerReadPhoto interface {
	Read() (models.Photos, error)
}

type HandlerForReadPhoto struct {
	log          logger.Logger
	photoActions PhotoActionsForHandlerReadPhoto
}

func NewHandlerForReadPhoto(log logger.Logger, photoActions PhotoActionsForHandlerReadPhoto) *HandlerForReadPhoto {
	result := &HandlerForReadPhoto{
		log:          log,
		photoActions: photoActions,
	}

	return result
}

func (handler *HandlerForReadPhoto) sendResponse(writer http.ResponseWriter, readPhotos models.Photos) {
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

func (handler *HandlerForReadPhoto) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	photos, err := handler.photoActions.Read()
	if err != nil {
		handler.log.Printf("can not read photos: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(writer, photos)
}
