package main

import (
	uploadfile "instagram/datasource/upload-file"
	handler_upload_photo "instagram/handlers/photos-files/upload"

	"instagram/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	log := logger.New()

	photosManager := uploadfile.New(log)

	handlerUploadPhoto := handler_upload_photo.NewHandlerForUploadPhoto(log, photosManager)
	router.Method(http.MethodPost, "/upload/{id_user}/{photo_name}", handlerUploadPhoto)
	//router.Method(http.MethodGet, "/download/{id_user}/{photo_name}", handlerDownloadPhoto)

	http.ListenAndServe(":8080", router)
}
