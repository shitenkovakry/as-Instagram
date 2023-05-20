package countoflikes

import (
	"encoding/json"
	"instagram/logger"
	"io"
	"net/http"

	models "instagram/models/likes"
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

func (handler *HandlerForCountLikes) prepareRequest(request *http.Request) (*models.Like, error) {
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

	var newCountOfLikes *CountLikes

	if err := json.Unmarshal(body, &newCountOfLikes); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	newCount := &models.Like{
		IdPhoto: newCountOfLikes.IDPhoto,
	}

	return newCount, nil
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
	newCount, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("cannot prepare request: %v", err)
		write.WriteHeader(http.StatusBadRequest)

		return
	}

	countedLikes, err := handler.likeActions.Count(newCount.IdPhoto)
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
