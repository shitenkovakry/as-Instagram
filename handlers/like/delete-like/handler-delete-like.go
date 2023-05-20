package deletelike

import (
	"encoding/json"
	"instagram/logger"
	"io"
	"net/http"

	models "instagram/models/likes"
)

type DeleteLike struct {
	IDUser  int `json:"user_id"`
	IDPhoto int `json:"photo_id"`
}

type LikeActionsForHandlerDeleteLike interface {
	Delete(photoID int, userID int) error
}

type HandlerForDeleteLike struct {
	log         logger.Logger
	likeActions LikeActionsForHandlerDeleteLike
}

func (handler *HandlerForDeleteLike) prepareRequest(request *http.Request) (*models.Like, error) {
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

	var deleteLikeFromClient *DeleteLike

	if err := json.Unmarshal(body, &deleteLikeFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	deletedLike := &models.Like{
		IDUser:  deleteLikeFromClient.IDUser,
		IDPhoto: deleteLikeFromClient.IDPhoto,
	}

	return deletedLike, nil
}

func (handler *HandlerForDeleteLike) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	IDs, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("can not prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := handler.likeActions.Delete(IDs.IDUser, IDs.IDPhoto); err != nil {
		handler.log.Printf("cannot create like: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func NewHandlerForDeleteLike(log logger.Logger, likeActions LikeActionsForHandlerDeleteLike) *HandlerForDeleteLike {
	result := &HandlerForDeleteLike{
		log:         log,
		likeActions: likeActions,
	}

	return result
}
