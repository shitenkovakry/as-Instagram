package delete

import (
	"encoding/json"
	"errors"
	"instagram/errordefs"
	"instagram/logger"
	models "instagram/models/comments"
	"io"
	"net/http"
)

type DeleteComment struct {
	IdComment int `json:"comment_id"`
}

type CommentActionsForHandlerDeleteComment interface {
	Delete(commentID int) (*models.Comment, error)
}

type HandlerForDeleteComment struct {
	log            logger.Logger
	commentActions CommentActionsForHandlerDeleteComment
}

func (handler *HandlerForDeleteComment) prepareRequest(request *http.Request) (*models.Comment, error) {
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

	var deleteCommentFromClient *DeleteComment

	if err := json.Unmarshal(body, &deleteCommentFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	deletedComment := &models.Comment{
		IdComment: deleteCommentFromClient.IdComment,
	}

	return deletedComment, nil
}

func (handler *HandlerForDeleteComment) sendResponse(writer http.ResponseWriter, deletedComment *models.Comment) {
	deletedCommentMarshaled, err := json.Marshal(deletedComment)
	if err != nil {
		handler.log.Printf("can not marshal deleted comment: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := writer.Write(deletedCommentMarshaled); err != nil {
		handler.log.Printf("can not send to client deleted comment: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForDeleteComment) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	shouldDeleteComment, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("can not prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	deletedComment, err := handler.commentActions.Delete(shouldDeleteComment.IdComment)
	if errors.Is(err, errordefs.ErrNoDocuments) {
		writer.WriteHeader(http.StatusNotFound)

		return
	}

	if err != nil {
		handler.log.Printf("can not delete comment: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(writer, deletedComment)
}

func NewHandlerForDeleteComment(log logger.Logger, commentActions CommentActionsForHandlerDeleteComment) *HandlerForDeleteComment {
	result := &HandlerForDeleteComment{
		log:            log,
		commentActions: commentActions,
	}

	return result
}
