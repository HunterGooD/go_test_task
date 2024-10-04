package entity

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal Server Error")
	// ErrNotFound will throw if the requested song is not exists
	ErrNotFound = errors.New("object is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("object already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("params is not valid")
)
