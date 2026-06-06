package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrAlreadyExists     = errors.New("already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrMappingExists     = errors.New("mapping already exists")
	ErrFrameworkInUse    = errors.New("framework has controls and cannot be deleted")
)
