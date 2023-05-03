package comment

import (
	"encoding/json"
	"io"
	"net/http"

	"instagram/logger"
	models "instagram/models/comments"
)

type AddComment struct {
	IDUser  int    `json:"user_id"`
	IDPhoto int    `json:"photo_id"`
	Comment string `json:"comment"`
}

type CommentActionsForHandlerAddComment interface {
	Add(commentOfUser *models.Comment) (*models.Comment, error)
}

type HandlerForAddComment struct {
	log            logger.Logger
	commentActions CommentActionsForHandlerAddComment
}

func (handler *HandlerForAddComment) prepareRequest(request *http.Request) (*models.Comment, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		handler.log.Printf("cannot read body: %v", err)

		return nil, err
	}

	var newCommentFromClient *AddComment

	if err := json.Unmarshal(body, &newCommentFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	newComment := &models.Comment{
		IdUser:  newCommentFromClient.IDUser,
		IdPhoto: newCommentFromClient.IDPhoto,
		Comment: newCommentFromClient.Comment,
	}

	return newComment, nil
}

func (handler *HandlerForAddComment) sendResponse(write http.ResponseWriter, createdComment *models.Comment) {
	createdCommentMarshaled, err := json.Marshal(createdComment)
	if err != nil {
		handler.log.Printf("cannot marshal created comment: %v", err)
		write.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := write.Write(createdCommentMarshaled); err != nil {
		handler.log.Printf("cannot send to client created comment: %v", err)
		write.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForAddComment) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	newComment, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("cannot prepare request: %v", err)
		write.WriteHeader(http.StatusBadRequest)

		return
	}

	createdComment, err := handler.commentActions.Add(newComment)
	if err != nil {
		handler.log.Printf("cannot create comment: %v", err)
		write.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(write, createdComment)
}

func NewHandlerForAddComment(log logger.Logger, commentActions CommentActionsForHandlerAddComment) *HandlerForAddComment {
	result := &HandlerForAddComment{
		log:            log,
		commentActions: commentActions,
	}

	return result
}
