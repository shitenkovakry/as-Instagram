package countoflikes

import (
	"encoding/json"
	"instagram/logger"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CountLikes struct {
	IDPhoto int `json:"photo_id"`
}

type LikeActionsForHandlerCountLikes interface {
	Count(idPhoto int) (int, error)
}

type HandlerForCountLikes struct {
	log         logger.Logger
	likeActions LikeActionsForHandlerCountLikes
}

func (handler *HandlerForCountLikes) prepareRequest(request *http.Request) (int, error) {
	photoIDParam := chi.URLParam(request, "id_photo")
	photoID, err := strconv.Atoi(photoIDParam)

	if err != nil {
		handler.log.Printf("err = %v", err)

		return 0, err
	}

	return photoID, nil
}

func (handler *HandlerForCountLikes) sendResponse(write http.ResponseWriter, countedLikes int) {
	createdCountedLikes, err := json.Marshal(countedLikes)
	if err != nil {
		handler.log.Printf("cannot marshal created counted likes: %v", err)
		write.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := write.Write(createdCountedLikes); err != nil {
		handler.log.Printf("cannot send to client counted likes: %v", err)
		write.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForCountLikes) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	photoID, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("cannot prepare request: %v", err)
		write.WriteHeader(http.StatusBadRequest)

		return
	}

	countedLikes, err := handler.likeActions.Count(photoID)
	if err != nil {
		handler.log.Printf("cannot count likes: %v", err)
		write.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(write, countedLikes)
}

func NewHandlerForCountLikes(log logger.Logger, likeActions LikeActionsForHandlerCountLikes) *HandlerForCountLikes {
	result := &HandlerForCountLikes{
		log:         log,
		likeActions: likeActions,
	}

	return result
}
