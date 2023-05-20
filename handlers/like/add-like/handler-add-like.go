package addlike

import (
	"encoding/json"
	"instagram/logger"
	"io"
	"net/http"

	models "instagram/models/likes"
)

type AddLike struct {
	IDUser  int `json:"user_id"`
	IDPhoto int `json:"photo_id"`
}

type LikeActionsForHandlerAddLike interface {
	Add(idPhoto int, idUser int) error
}

type HandlerForAddLike struct {
	log         logger.Logger
	likeActions LikeActionsForHandlerAddLike
}

func (handler *HandlerForAddLike) prepareRequest(request *http.Request) (*models.Like, error) {
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

	var newLikeFromClient *AddLike

	if err := json.Unmarshal(body, &newLikeFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	newLike := &models.Like{
		IDUser:  newLikeFromClient.IDUser,
		IDPhoto: newLikeFromClient.IDPhoto,
	}

	return newLike, nil
}

func (handler *HandlerForAddLike) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	newLike, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("cannot prepare request: %v", err)
		write.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := handler.likeActions.Add(newLike.IDPhoto, newLike.IDUser); err != nil {
		handler.log.Printf("cannot create like: %v", err)
		write.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func NewHandlerForAddLike(log logger.Logger, likeActions LikeActionsForHandlerAddLike) *HandlerForAddLike {
	result := &HandlerForAddLike{
		log:         log,
		likeActions: likeActions,
	}

	return result
}
