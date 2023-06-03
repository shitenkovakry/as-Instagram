package readcomments

import (
	"encoding/json"
	"instagram/logger"
	models "instagram/models/comments"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type ReadComments struct {
	IdPhoto int `json:"id_photo"`
	IdUser  int `json:"id_user"`
}

type CommentsActionsForHandlerReadComments interface {
	ReadComments(userID int, photoID int) (models.Comments, error)
}

type HandlerForReadComments struct {
	log             logger.Logger
	commentsActions CommentsActionsForHandlerReadComments
}

func NewHandlerForReadComments(log logger.Logger, commentsActions CommentsActionsForHandlerReadComments) *HandlerForReadComments {
	result := &HandlerForReadComments{
		log:             log,
		commentsActions: commentsActions,
	}

	return result
}

func (handler *HandlerForReadComments) prepareRequest(request *http.Request) (*models.Comment, error) {
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

	var readCommentFromClient *ReadComments

	if err := json.Unmarshal(body, &readCommentFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	readComm := &models.Comment{
		IdUser:  userID,
		IdPhoto: photoID,
	}

	return readComm, nil
}

func (handler *HandlerForReadComments) sendResponse(writer http.ResponseWriter, readComments models.Comments) {
	data, err := json.Marshal(readComments)
	if err != nil {
		handler.log.Printf("can not marshal comments: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := writer.Write(data); err != nil {
		handler.log.Printf("can not writer comments: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForReadComments) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	shouldReadComments, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("can not prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	comments, err := handler.commentsActions.ReadComments(shouldReadComments.IdUser, shouldReadComments.IdPhoto)
	if err != nil {
		handler.log.Printf("can not read comments: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(writer, comments)
}
