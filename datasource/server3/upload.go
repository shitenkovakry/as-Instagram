package server3

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"instagram/logger"
)

const (
	basePath = "./photos"
)

type FileManager struct {
	log logger.Logger
}

func (fileManager *FileManager) Upload(userID int, photoContent []byte, photoFilename string) error {
	path := fmt.Sprintf("%s/%d/%s", basePath, userID, photoFilename)

	if err := os.WriteFile(path, photoContent, os.ModePerm); err != nil {
		return errors.Wrap(err, "cannot save photo content")
	}

	return nil
}

func New(log logger.Logger) *FileManager {
	return &FileManager{
		log: log,
	}
}
