package domain

import (
	"time"

	"github.com/google/uuid"
)

type Framework struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	Description string
	Version     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewFramework(name, slug, description, version string) *Framework {
	return &Framework{
		ID:          uuid.New(),
		Name:        name,
		Slug:        slug,
		Description: description,
		Version:     version,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

type ControlCategory string

const (
	ControlCategoryAdministrative ControlCategory = "administrative"
	ControlCategoryTechnical      ControlCategory = "technical"
	ControlCategoryPhysical       ControlCategory = "physical"
)

type Control struct {
	ID              uuid.UUID
	FrameworkID     uuid.UUID
	ControlID       string
	Name            string
	Description     string
	Category        ControlCategory
	RiskDescription string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewControl(frameworkID uuid.UUID, controlID, name, description string, category ControlCategory, riskDesc string) *Control {
	return &Control{
		ID:              uuid.New(),
		FrameworkID:     frameworkID,
		ControlID:       controlID,
		Name:            name,
		Description:     description,
		Category:        category,
		RiskDescription: riskDesc,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

type MappingStatus string

const (
	MappingStatusMapped    MappingStatus = "mapped"
	MappingStatusPartial   MappingStatus = "partial"
	MappingStatusUnmapped  MappingStatus = "unmapped"
)

type EvidenceMapping struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	EvidenceID  uuid.UUID
	ControlID   uuid.UUID
	Notes       string
	MappedBy    uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewEvidenceMapping(tenantID, evidenceID, controlID, mappedBy uuid.UUID, notes string) *EvidenceMapping {
	return &EvidenceMapping{
		ID:         uuid.New(),
		TenantID:   tenantID,
		EvidenceID: evidenceID,
		ControlID:  controlID,
		Notes:      notes,
		MappedBy:   mappedBy,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

type QuestionMapping struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	QuestionID   uuid.UUID
	ControlID    uuid.UUID
	MappedBy     uuid.UUID
	CreatedAt    time.Time
}

type ControlCoverage struct {
	Control     Control
	Status      MappingStatus
	EvidenceIDs []uuid.UUID
}

type FrameworkCoverage struct {
	Framework Framework
	Controls  []ControlCoverage
	Total     int
	Mapped    int
}
