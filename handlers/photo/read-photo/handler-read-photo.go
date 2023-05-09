package readphoto

import (
	"encoding/json"
	"instagram/logger"
	models "instagram/models/photos"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ReadPhoto struct {
	IDUser  int `json:"id_user"`
	IDPhoto int `json:"id_photo"`
}

type PhotoActionsForHandlerReadPhoto interface {
	ReadPhoto(idUser int, idPhoto int) (*models.Photo, error)
}

type HandlerForReadPhoto struct {
	log          logger.Logger
	photoActions PhotoActionsForHandlerReadPhoto
}

func (handler *HandlerForReadPhoto) prepareRequest(request *http.Request) (*models.Photo, error) {
	defer func() {
		if err := request.Body.Close(); err != nil {
			handler.log.Printf("cannot close body: %v", err)
		}
	}()

	userIDParam := chi.URLParam(request, "id_user")
	userID, err := strconv.Atoi(userIDParam)

	if err != nil {
		handler.log.Printf("err = %v", err)

		return nil, err
	}

	photoIDParam := chi.URLParam(request, "id_photo")
	photoID, err := strconv.Atoi(photoIDParam)

	if err != nil {
		handler.log.Printf("err = %v", err)

		return nil, err
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		handler.log.Printf("cannot read body: %v", err)

		return nil, err
	}

	var readPhotoFromClient *ReadPhoto

	if err := json.Unmarshal(body, &readPhotoFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	readPhoto := &models.Photo{
		IDUser:  userID,
		IDPhoto: photoID,
	}

	return readPhoto, nil
}

func (handler *HandlerForReadPhoto) sendResponse(writer http.ResponseWriter, readPhoto *models.Photo) {
	data, err := json.Marshal(readPhoto)
	if err != nil {
		handler.log.Printf("can not marshal photo: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := writer.Write(data); err != nil {
		handler.log.Printf("can not writer photo: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForReadPhoto) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	shouldReadPhoto, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("can not prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	photo, err := handler.photoActions.ReadPhoto(shouldReadPhoto.IDUser, shouldReadPhoto.IDPhoto)
	if err != nil {
		handler.log.Printf("can not read photos: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(writer, photo)
}

func NewHandlerForReadPhoto(log logger.Logger, photoActions PhotoActionsForHandlerReadPhoto) *HandlerForReadPhoto {
	result := &HandlerForReadPhoto{
		log:          log,
		photoActions: photoActions,
	}

	return result
}
