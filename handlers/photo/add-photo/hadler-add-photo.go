package addphoto

import (
	"encoding/json"
	"instagram/errordefs"
	"instagram/logger"
	models "instagram/models/photos"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

const maxMem = (2 << 21) * 2

type PhotoActionsForHandlerAddPhoto interface {
	Add(userID int, photoContent []byte, photoFilename string) (*models.Photo, error)
}

type HandlerForAddPhoto struct {
	log          logger.Logger
	photoActions PhotoActionsForHandlerAddPhoto
}

type PhotoRequest struct {
	IDUser   int
	Content  multipart.File
	Filename string
}

func (handler *HandlerForAddPhoto) prepareRequest(request *http.Request) (*PhotoRequest, error) {
	if err := request.ParseMultipartForm(maxMem); err != nil {
		return nil, errors.Wrapf(err, "can not parse to multipart")
	} // - читает информацию requesta, а request содержит body и эта функция просматривает  информацию body: сколько весит бинарный и тд

	userIDs, ok := request.MultipartForm.Value["user-id"]
	if !ok || len(userIDs) == 0 {
		return nil, errors.Wrapf(errordefs.ErrNotFound, "can not read the user-id")
	} // из этой функции считываем юсер айди

	userID, err := strconv.Atoi(userIDs[0])
	if err != nil {
		return nil, errors.Wrapf(err, "can not convert")
	} // айди из стринга переводится в инт

	filesUpload, ok := request.MultipartForm.File["file-upload"]
	if !ok || len(filesUpload) == 0 {
		return nil, errors.Wrapf(err, "can not send at least one file")
	} // считываем файл аплоад(нашу фотку)

	fileUpload := filesUpload[0]
	file, err := fileUpload.Open()
	if err != nil {
		return nil, errors.Wrapf(errordefs.ErrNotFound, "can not read the file")
	}

	newPhoto := &PhotoRequest{
		IDUser:   userID,
		Filename: fileUpload.Filename,
		Content:  file, //само фото в виде файла, а файл содержит бинарный формат фото
	}

	return newPhoto, nil
}

func (handler *HandlerForAddPhoto) sendResponse(writes http.ResponseWriter, addedPhoto *models.Photo) {
	addedPhotoMarshaled, err := json.Marshal(addedPhoto)
	if err != nil {
		handler.log.Printf("cannot marshal added photo: %v", err)
		writes.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := writes.Write(addedPhotoMarshaled); err != nil {
		handler.log.Printf("cannot send to client added photo: %v", err)
		writes.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (handler *HandlerForAddPhoto) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	newPhoto, err := handler.prepareRequest(request)
	if err != nil {
		handler.log.Printf("cannot prepare request: %v", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	content, err := io.ReadAll(newPhoto.Content)
	if err != nil {
		panic(err)
	}

	addedPhoto, err := handler.photoActions.Add(newPhoto.IDUser, content, newPhoto.Filename)
	if err != nil {
		handler.log.Printf("cannot add photo: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.sendResponse(writer, addedPhoto)
}

func NewHandlerForAddPhoto(log logger.Logger, photoActions PhotoActionsForHandlerAddPhoto) *HandlerForAddPhoto {
	result := &HandlerForAddPhoto{
		log:          log,
		photoActions: photoActions,
	}

	return result
}
