package domain

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrInvalidFormat    = errors.New("invalid export format")
	ErrExportFailed     = errors.New("export generation failed")
	ErrEvidenceNotFound = errors.New("evidence not found")
)
