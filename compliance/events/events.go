package events

import (
	"github.com/google/uuid"
)

const (
	SubjectFrameworkCreated    = "compliance.framework.created"
	SubjectControlCreated      = "compliance.control.created"
	SubjectEvidenceMapped      = "compliance.evidence.mapped"
	SubjectEvidenceUnmapped    = "compliance.evidence.unmapped"
)

type FrameworkCreated struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

func (e FrameworkCreated) Subject() string { return SubjectFrameworkCreated }

type ControlCreated struct {
	ID          uuid.UUID `json:"id"`
	FrameworkID uuid.UUID `json:"framework_id"`
	ControlID   string    `json:"control_id"`
	Name        string    `json:"name"`
}

func (e ControlCreated) Subject() string { return SubjectControlCreated }

type EvidenceMapped struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	EvidenceID uuid.UUID `json:"evidence_id"`
	ControlID  uuid.UUID `json:"control_id"`
}

func (e EvidenceMapped) Subject() string { return SubjectEvidenceMapped }

type EvidenceUnmapped struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	EvidenceID uuid.UUID `json:"evidence_id"`
	ControlID  uuid.UUID `json:"control_id"`
}

func (e EvidenceUnmapped) Subject() string { return SubjectEvidenceUnmapped }
