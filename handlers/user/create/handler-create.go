package create

import (
	"encoding/json"
	"io"
	"net/http"

	"instagram/logger"
	models "instagram/models/users"
)

type NewUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserActionsForHandlerCreateUser interface {
	Create(NewUser *models.UserRegistration) (*models.UserRegistration, error)
}

type HandlerForCreateUser struct {
	log         logger.Logger
	userActions UserActionsForHandlerCreateUser
}

func NewHandlerForCreateUser(log logger.Logger, userActions UserActionsForHandlerCreateUser) *HandlerForCreateUser {
	result := &HandlerForCreateUser{
		log:         log,
		userActions: userActions,
	}

	return result
}

func (handler *HandlerForCreateUser) prepareRequest(request *http.Request) (*models.UserRegistration, error) {
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

	var newUserFromClient *NewUser

	if err := json.Unmarshal(body, &newUserFromClient); err != nil {
		handler.log.Printf("cannot unmarshal body=%s: %v", string(body), err)

		return nil, err
	}

	newUser := &models.UserRegistration{
		Name:  newUserFromClient.Name,
		Email: newUserFromClient.Email,
	}

	return newUser, nil
}

func (handler *HandlerForCreateUser) sendResponse(write http.ResponseWriter, createdUser *models.UserRegistration) {
	createdUserMarshaled, err := json.Marshal(createdUser)
	if err != nil {
		handler.log.Printf("cannot marshal created user: %v", err)
		write.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := write.Write(createdUserMarshaled); err != nil {
		handler.log.Printf("cannot send to client created user: %v", err)
		write.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForCreateUser) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	handler.log.Printf("create request hit")

	newUser, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("cannot prepare request: %v", err)
		write.WriteHeader(http.StatusBadRequest)

		return
	}

	createdUser, err := handler.userActions.Create(newUser)
	if err != nil {
		handler.log.Printf("cannot create user: %v", err)
		write.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(write, createdUser)
}
