package readuser

import (
	"encoding/json"
	"instagram/logger"
	models "instagram/models/users"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type ReadUser struct {
	ID int `json:"user_id"`
}

type UserActionsForHandlerReadUser interface {
	ReadUser(idUser int) (*models.UserRegistration, error)
}

type HandlerForReadUser struct {
	log         logger.Logger
	userActions UserActionsForHandlerReadUser
}

func (handler *HandlerForReadUser) prepareRequest(request *http.Request) (*models.UserRegistration, error) {
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

	readUser := &models.UserRegistration{
		ID: userID,
	}

	return readUser, nil
}

func (handler *HandlerForReadUser) sendResponse(writer http.ResponseWriter, readUser *models.UserRegistration) {
	data, err := json.Marshal(readUser)
	if err != nil {
		handler.log.Printf("can not marshal user: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := writer.Write(data); err != nil {
		handler.log.Printf("can not writer user: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForReadUser) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	shouldReadUser, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("can not prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	user, err := handler.userActions.ReadUser(shouldReadUser.ID)
	if err != nil {
		handler.log.Printf("can not read user: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(writer, user)
}

func NewHandlerForReadUser(log logger.Logger, userActions UserActionsForHandlerReadUser) *HandlerForReadUser {
	result := &HandlerForReadUser{
		log:         log,
		userActions: userActions,
	}

	return result
}
