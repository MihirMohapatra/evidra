package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrInvalidStatus    = errors.New("invalid status transition")
	ErrUnsupportedFile  = errors.New("unsupported file type")
	ErrParseFailed      = errors.New("document parse failed")
	ErrFileTooLarge     = errors.New("file too large")
)

func NewStatusError(current, target Status) error {
	return fmt.Errorf("%w: cannot transition from %s to %s", ErrInvalidStatus, current, target)
}
