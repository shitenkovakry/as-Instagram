package main

import (
	server3 "instagram/datasource/server3"
	handler_download_photo "instagram/handlers/photos-files/download"
	handler_upload_photo "instagram/handlers/photos-files/upload"

	"instagram/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	log := logger.New()

	photosManager := server3.New(log)

	handlerUploadPhoto := handler_upload_photo.NewHandlerForUploadPhoto(log, photosManager)
	router.Method(http.MethodPost, "/upload/{id_user}/{photo_name}", handlerUploadPhoto)

	handlerDownloadPhoto := handler_download_photo.NewHandlerForDownloadPhoto(log, photosManager)
	router.Method(http.MethodGet, "/download/{id_user}/{photo_name}", handlerDownloadPhoto)

	http.ListenAndServe(":8080", router)
}
