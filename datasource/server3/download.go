package server3

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func (fileManager *FileManager) Download(photoFilename string, userID int) ([]byte, error) {
	path := fmt.Sprintf("%s/%d/%s", basePath, userID, photoFilename)

	photoContent, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "cannot download photo content")
	}

	return photoContent, nil
}
