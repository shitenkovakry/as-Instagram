package update

import (
	"encoding/json"
	"instagram/logger"
	models "instagram/models/users"
	"io"
	"net/http"
)

type UpdateUserByEmail struct {
	Id    int    `json:"user_id"`
	Email string `json:"new email"`
}

type UserActionsForHandlerUpdateUserByEmail interface {
	UpdateByEmail(userID int, newEmail string) (*models.UserRegistration, error)
}

type HandlerForUpdateEmail struct {
	log         logger.Logger
	userActions UserActionsForHandlerUpdateUserByEmail
}

func (handler *HandlerForUpdateEmail) prepareRequest(request *http.Request) (*models.UserRegistration, error) {
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

	var newNameFromClient *UpdateUserByEmail

	if err := json.Unmarshal(body, &newNameFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	newName := &models.UserRegistration{
		ID:    newNameFromClient.Id,
		Email: newNameFromClient.Email,
	}

	return newName, nil
}

func (handler *HandlerForUpdateEmail) sendResponse(writer http.ResponseWriter, updatedEmail *models.UserRegistration) {
	updatedEmailMarshaled, err := json.Marshal(updatedEmail)
	if err != nil {
		handler.log.Printf("cannot marshal updated email: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := writer.Write(updatedEmailMarshaled); err != nil {
		handler.log.Printf("cannot send to client updated email: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForUpdateEmail) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	shouldUpdateEmail, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("cannot prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	updatedEmail, err := handler.userActions.UpdateByEmail(shouldUpdateEmail.ID, shouldUpdateEmail.Email)
	if err != nil {
		handler.log.Printf("cannot update email: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(writer, updatedEmail)
}

func NewHandlerForUpdateUserByEmail(log logger.Logger, userActions UserActionsForHandlerUpdateUserByEmail) *HandlerForUpdateEmail {
	result := &HandlerForUpdateEmail{
		log:         log,
		userActions: userActions,
	}

	return result
}
