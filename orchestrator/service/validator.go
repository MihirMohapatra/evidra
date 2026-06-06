package service

import (
	"fmt"

	"github.com/evidra/evidra/orchestrator/domain"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateDraft(draft *domain.Draft) error {
	if draft.Answer == "" {
		return fmt.Errorf("%w: answer cannot be empty", domain.ErrValidationFailed)
	}
	if draft.Confidence < 0 || draft.Confidence > 1 {
		return fmt.Errorf("%w: confidence must be between 0 and 1", domain.ErrValidationFailed)
	}
	if draft.QuestionText == "" {
		return fmt.Errorf("%w: question text cannot be empty", domain.ErrValidationFailed)
	}
	if draft.ModelUsed == "" {
		return fmt.Errorf("%w: model used cannot be empty", domain.ErrValidationFailed)
	}
	return nil
}
