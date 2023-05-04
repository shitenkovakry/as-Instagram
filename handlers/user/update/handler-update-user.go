package update

import (
	"encoding/json"
	"instagram/logger"
	models "instagram/models/users"
	"io"
	"net/http"
)

type UpdateUserByName struct {
	Id   int    `json:"user_id"`
	Name string `json:"new name"`
}

type UserActionsForHandlerUpdateUserByName interface {
	Update(userID int, newName string) (*models.UserRegistration, error)
}

type HandlerForUpdateName struct {
	log         logger.Logger
	userActions UserActionsForHandlerUpdateUserByName
}

func (handler *HandlerForUpdateName) prepareRequest(request *http.Request) (*models.UserRegistration, error) {
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

	var newNameFromClient *UpdateUserByName

	if err := json.Unmarshal(body, &newNameFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	newName := &models.UserRegistration{
		ID:   newNameFromClient.Id,
		Name: newNameFromClient.Name,
	}

	return newName, nil
}

func (handler *HandlerForUpdateName) sendResponse(writer http.ResponseWriter, updatedName *models.UserRegistration) {
	updatedNameMarshaled, err := json.Marshal(updatedName)
	if err != nil {
		handler.log.Printf("cannot marshal updated name: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := writer.Write(updatedNameMarshaled); err != nil {
		handler.log.Printf("cannot send to client updated name: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForUpdateName) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	shouldUpdateName, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("cannot prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	updatedName, err := handler.userActions.Update(shouldUpdateName.ID, shouldUpdateName.Name)
	if err != nil {
		handler.log.Printf("cannot update name: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(writer, updatedName)
}

func NewHandlerForUpdateUserByName(log logger.Logger, userActions UserActionsForHandlerUpdateUserByName) *HandlerForUpdateName {
	result := &HandlerForUpdateName{
		log:         log,
		userActions: userActions,
	}

	return result
}
