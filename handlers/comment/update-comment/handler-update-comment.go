package updatecomment

import (
	"encoding/json"
	"instagram/logger"
	models "instagram/models/comments"
	"io"
	"net/http"
)

type UpdateComment struct {
	IdComment int    `json:"comment_id"`
	Comment   string `json:"new comment"`
}

type CommentActionsForHandlerUpdateComment interface {
	Update(commentID int, newComment string) (*models.Comment, error)
}

type HandlerForUpdateComment struct {
	log            logger.Logger
	commentActions CommentActionsForHandlerUpdateComment
}

func (handler *HandlerForUpdateComment) prepareRequest(request *http.Request) (*models.Comment, error) {
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

	var newCommentFromClient *UpdateComment

	if err := json.Unmarshal(body, &newCommentFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	newComment := &models.Comment{
		IdComment: newCommentFromClient.IdComment,
		Comment:   newCommentFromClient.Comment,
	}

	return newComment, nil
}

func (handler *HandlerForUpdateComment) sendResponse(writer http.ResponseWriter, updatedComment *models.Comment) {
	updatedCommentMarshaled, err := json.Marshal(updatedComment)
	if err != nil {
		handler.log.Printf("cannot marshal updated comment: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := writer.Write(updatedCommentMarshaled); err != nil {
		handler.log.Printf("cannot send to client updated comment: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForUpdateComment) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	shouldUpdateComment, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("cannot prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	updatedComment, err := handler.commentActions.Update(shouldUpdateComment.IdComment, shouldUpdateComment.Comment)
	if err != nil {
		handler.log.Printf("cannot update comment: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(writer, updatedComment)
}

func NewHandlerForUpdateComment(log logger.Logger, commentActions CommentActionsForHandlerUpdateComment) *HandlerForUpdateComment {
	result := &HandlerForUpdateComment{
		log:            log,
		commentActions: commentActions,
	}

	return result
}
