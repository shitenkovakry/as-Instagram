package upload

import (
	"io"
	"net/http"
	"strconv"

	"instagram/logger"

	"github.com/go-chi/chi"
)

type PhotoActionsForHandlerUploadPhoto interface {
	Upload(userID int, photoContent []byte, photoFilename string) error
}

type HandlerForUploadPhoto struct {
	log          logger.Logger
	photoActions PhotoActionsForHandlerUploadPhoto
}

func (handler *HandlerForUploadPhoto) prepareRequest(request *http.Request) (int, []byte, string, error) {
	defer func() {
		if err := request.Body.Close(); err != nil {
			handler.log.Printf("cannot close body: %v", err)
		}
	}()

	photoContent, err := io.ReadAll(request.Body)
	if err != nil {
		handler.log.Printf("cannot read photoContent: %v", err)

		return 0, nil, "", err
	}

	userIDParam := chi.URLParam(request, "id_user")
	userID, err := strconv.Atoi(userIDParam)

	if err != nil {
		handler.log.Printf("err = %v", err)

		return 0, nil, "", err
	}

	photoName := chi.URLParam(request, "photo_name")
	if photoName == "" {
		handler.log.Printf("can not accept empty photo name")

		return 0, nil, "", err
	}

	return userID, photoContent, photoName, nil
}

func (handler *HandlerForUploadPhoto) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	userID, photoContent, photoFilename, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("cannot prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := handler.photoActions.Upload(userID, photoContent, photoFilename); err != nil {
		handler.log.Printf("cannot upload photo: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func NewHandlerForUploadPhoto(log logger.Logger, photoActions PhotoActionsForHandlerUploadPhoto) *HandlerForUploadPhoto {
	result := &HandlerForUploadPhoto{
		log:          log,
		photoActions: photoActions,
	}

	return result
}
