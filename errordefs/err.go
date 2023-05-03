package errordefs

import "errors"

var (
	ErrNoDocuments = errors.New("no documents")
	ErrNotFound    = errors.New("not found")
	ErrAlready     = errors.New("already")
)
