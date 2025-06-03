package storage

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	errAppNotFound       = errors.New("app not found")
)
