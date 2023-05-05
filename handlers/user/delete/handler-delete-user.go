package delete

import (
	"encoding/json"
	"errors"
	"instagram/errordefs"
	"instagram/logger"
	models "instagram/models/users"
	"io"
	"net/http"
)

type DeleteUser struct {
	Id int `json:"user_id"`
}

type UserActionsForHandlerDeleteUser interface {
	Delete(userID int) (*models.UserRegistration, error)
}

type HandlerForDeleteUser struct {
	log         logger.Logger
	userActions UserActionsForHandlerDeleteUser
}

func (handler *HandlerForDeleteUser) prepareRequest(request *http.Request) (*models.UserRegistration, error) {
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

	var deleteUserFromClient *DeleteUser

	if err := json.Unmarshal(body, &deleteUserFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	deletedUser := &models.UserRegistration{
		ID: deleteUserFromClient.Id,
	}

	return deletedUser, nil
}

func (handler *HandlerForDeleteUser) sendResponse(writer http.ResponseWriter, deletedUser *models.UserRegistration) {
	deletedUserMarshaled, err := json.Marshal(deletedUser)
	if err != nil {
		handler.log.Printf("can not marshal deleted user: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := writer.Write(deletedUserMarshaled); err != nil {
		handler.log.Printf("can not send to client deleted user: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForDeleteUser) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	shouldDeleteUser, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("can not prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	deletedUser, err := handler.userActions.Delete(shouldDeleteUser.ID)
	if errors.Is(err, errordefs.ErrNoDocuments) {
		writer.WriteHeader(http.StatusNotFound)

		return
	}

	if err != nil {
		handler.log.Printf("can not delete user: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(writer, deletedUser)
}

func NewHandlerForDeleteUser(log logger.Logger, userActions UserActionsForHandlerDeleteUser) *HandlerForDeleteUser {
	result := &HandlerForDeleteUser{
		log:         log,
		userActions: userActions,
	}

	return result
}
