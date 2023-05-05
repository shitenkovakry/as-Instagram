package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"instagram/datasource/managers/comments"
	photos "instagram/datasource/managers/photos"
	users "instagram/datasource/managers/users"
	"instagram/datasource/mongo"
	handler_add "instagram/handlers/comment/add-comment"
	handler_delete_comment "instagram/handlers/comment/delete-comment"
	handler_update_comment "instagram/handlers/comment/update-comment"
	handler_add_photo "instagram/handlers/photo/add-photo"
	handler_delete_photo "instagram/handlers/photo/delete-photo"
	handler_read "instagram/handlers/photo/read-photo"
	handler_create "instagram/handlers/user/create"
	handler_delete_user "instagram/handlers/user/delete"
	handler_update_email "instagram/handlers/user/update"
	handler_update_name "instagram/handlers/user/update"
	"instagram/logger"

	"github.com/go-chi/chi/v5"
)

const (
	addr = ":8080"
)

// docker run -d --rm --name my_mongo -p 27017:27017 mongo:latest

func main() {
	router := chi.NewRouter()
	log := logger.New()

	usersDB := mongo.NewUsersManager(log, "", "", []string{"localhost:27017"}, "my-database")
	photosDB := mongo.NewPhotosManager(log, "", "", []string{"localhost:27017"}, "my-database")
	commentsDB := mongo.NewCommentsManager(log, "", "", []string{"localhost:27017"}, "my-database")

	usersManager := users.New(log, usersDB)
	photosManager := photos.New(log, photosDB)
	commentsManager := comments.New(log, commentsDB)

	handlerForCreateUser := handler_create.NewHandlerForCreateUser(log, usersManager)
	router.Method(http.MethodPost, "/api/v1/users/register", handlerForCreateUser)
	handlerForUpdateNameOfUser := handler_update_name.NewHandlerForUpdateUserByName(log, usersManager)
	router.Method(http.MethodPut, "/api/v1/users/update/name", handlerForUpdateNameOfUser)
	handlerForUpdateEmailOfUser := handler_update_email.NewHandlerForUpdateUserByEmail(log, usersManager)
	router.Method(http.MethodPut, "/api/v1/users/update/email", handlerForUpdateEmailOfUser)
	handlerDeleteUser := handler_delete_user.NewHandlerForDeleteUser(log, usersManager)
	router.Method(http.MethodDelete, "/api/v1/users/delete", handlerDeleteUser)

	handlerReadPhotos := handler_read.NewHandlerForReadPhoto(log, photosManager)
	router.Method(http.MethodGet, "/api/v1/photos", handlerReadPhotos)
	handlerAddPhoto := handler_add_photo.NewHandlerForAddPhoto(log, photosManager)
	router.Method(http.MethodPost, "/api/v1/photos/add", handlerAddPhoto)
	handlerDeletePhoto := handler_delete_photo.NewHandlerForDeletePhoto(log, photosManager)
	router.Method(http.MethodDelete, "/api/v1/photos/delete", handlerDeletePhoto)

	handlerAddComment := handler_add.NewHandlerForAddComment(log, commentsManager)
	router.Method(http.MethodPost, "/api/v1/comments/add", handlerAddComment)
	handlerDeleteComment := handler_delete_comment.NewHandlerForDeleteComment(log, commentsManager)
	router.Method(http.MethodDelete, "/api/v1/comments/delete", handlerDeleteComment)
	handlerForUpdateComment := handler_update_comment.NewHandlerForUpdateComment(log, commentsManager)
	router.Method(http.MethodPut, "/api/v1/comments/update", handlerForUpdateComment)

	server := NewServer(addr, router)

	log.Printf("Serving at [%s]", addr)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server is error: %v", err)
		}
	}()

	<-stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// завершение работы серверов
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	if err := usersDB.Shutdown(ctx); err != nil {
		log.Printf("userDB error: %v", err)
	}

	if err := photosDB.Shutdown(ctx); err != nil {
		log.Printf("photosDB error: %v", err)
	}

	if err := commentsDB.Shutdown(ctx); err != nil {
		log.Printf("commentsDB error: %v", err)
	}
}

func NewServer(address string, router *chi.Mux) *http.Server {
	return &http.Server{
		Addr:    address,
		Handler: router,
	}
}
