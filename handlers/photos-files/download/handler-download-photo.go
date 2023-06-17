package download

import (
	"instagram/logger"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type PhotoActionsForHandlerDownloadPhoto interface {
	Download(photoFilename string, userID int) ([]byte, error)
}

type HandlerForDownloadPhoto struct {
	log          logger.Logger
	photoActions PhotoActionsForHandlerDownloadPhoto
}

func (handler *HandlerForDownloadPhoto) prepareRequest(request *http.Request) (string, int, error) {
	defer func() {
		if err := request.Body.Close(); err != nil {
			handler.log.Printf("cannot close body: %v", err)
		}
	}()

	userIDParam := chi.URLParam(request, "id_user")
	userID, err := strconv.Atoi(userIDParam)

	if err != nil {
		handler.log.Printf("err = %v", err)

		return "", 0, nil
	}

	photoName := chi.URLParam(request, "photo_name")
	if photoName == "" {
		handler.log.Printf("can not accept empty photo name")

		return "", 0, nil
	}

	return photoName, userID, nil
}

func (handler *HandlerForDownloadPhoto) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	photoFilename, userID, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("cannot prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	content, err := handler.photoActions.Download(photoFilename, userID)
	if err != nil {
		handler.log.Printf("cannot upload photo: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := writer.Write(content); err != nil {
		handler.log.Printf("cannot send to client downloaded photo: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func NewHandlerForDownloadPhoto(log logger.Logger, photoActions PhotoActionsForHandlerDownloadPhoto) *HandlerForDownloadPhoto {
	result := &HandlerForDownloadPhoto{
		log:          log,
		photoActions: photoActions,
	}

	return result
}
