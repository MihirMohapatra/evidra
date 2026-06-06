package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrAlreadyExists     = errors.New("already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrInvalidCategory   = errors.New("invalid category")
	ErrEvidenceExpired   = errors.New("evidence has expired")
	ErrNotOwner          = errors.New("user is not the owner")
)
