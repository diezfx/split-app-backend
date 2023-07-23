package storage

import "errors"

var (
	ErrNotFound      = errors.New("element not found")
	ErrAlreadyExists = errors.New("element already exists")
)
