package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/compliance/domain"
)

// --- FrameworkRepo ---

type FrameworkRepo struct {
	pool *pgxpool.Pool
}

func NewFrameworkRepo(pool *pgxpool.Pool) *FrameworkRepo {
	return &FrameworkRepo{pool: pool}
}

func (r *FrameworkRepo) Create(ctx context.Context, f *domain.Framework) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO compliance_frameworks (id, name, slug, description, version, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		f.ID, f.Name, f.Slug, f.Description, f.Version, f.CreatedAt, f.UpdatedAt)
	return err
}

func (r *FrameworkRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Framework, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, description, version, created_at, updated_at
		 FROM compliance_frameworks WHERE id = $1`, id)
	return scanFramework(row)
}

func (r *FrameworkRepo) GetBySlug(ctx context.Context, slug string) (*domain.Framework, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, description, version, created_at, updated_at
		 FROM compliance_frameworks WHERE slug = $1`, slug)
	return scanFramework(row)
}

func (r *FrameworkRepo) List(ctx context.Context) ([]*domain.Framework, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, slug, description, version, created_at, updated_at
		 FROM compliance_frameworks ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFrameworks(rows)
}

func (r *FrameworkRepo) Update(ctx context.Context, f *domain.Framework) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE compliance_frameworks SET name = $1, slug = $2, description = $3, version = $4, updated_at = $5 WHERE id = $6`,
		f.Name, f.Slug, f.Description, f.Version, f.UpdatedAt, f.ID)
	return err
}

func (r *FrameworkRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM compliance_frameworks WHERE id = $1`, id)
	return err
}

// --- ControlRepo ---

type ControlRepo struct {
	pool *pgxpool.Pool
}

func NewControlRepo(pool *pgxpool.Pool) *ControlRepo {
	return &ControlRepo{pool: pool}
}

func (r *ControlRepo) Create(ctx context.Context, c *domain.Control) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO compliance_controls (id, framework_id, control_id, name, description, category, risk_description, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		c.ID, c.FrameworkID, c.ControlID, c.Name, c.Description, c.Category, c.RiskDescription, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *ControlRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Control, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, framework_id, control_id, name, description, category, risk_description, created_at, updated_at
		 FROM compliance_controls WHERE id = $1`, id)
	return scanControl(row)
}

func (r *ControlRepo) ListByFramework(ctx context.Context, frameworkID uuid.UUID) ([]*domain.Control, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, framework_id, control_id, name, description, category, risk_description, created_at, updated_at
		 FROM compliance_controls WHERE framework_id = $1 ORDER BY control_id`, frameworkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanControls(rows)
}

func (r *ControlRepo) Update(ctx context.Context, c *domain.Control) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE compliance_controls SET name = $1, description = $2, category = $3, risk_description = $4, updated_at = $5 WHERE id = $6`,
		c.Name, c.Description, c.Category, c.RiskDescription, c.UpdatedAt, c.ID)
	return err
}

func (r *ControlRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM compliance_controls WHERE id = $1`, id)
	return err
}

func (r *ControlRepo) CountByFramework(ctx context.Context, frameworkID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM compliance_controls WHERE framework_id = $1`, frameworkID).Scan(&count)
	return count, err
}

// --- EvidenceMappingRepo ---

type EvidenceMappingRepo struct {
	pool *pgxpool.Pool
}

func NewEvidenceMappingRepo(pool *pgxpool.Pool) *EvidenceMappingRepo {
	return &EvidenceMappingRepo{pool: pool}
}

func (r *EvidenceMappingRepo) Create(ctx context.Context, m *domain.EvidenceMapping) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO compliance_evidence_mappings (id, tenant_id, evidence_id, control_id, notes, mapped_by, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		m.ID, m.TenantID, m.EvidenceID, m.ControlID, m.Notes, m.MappedBy, m.CreatedAt, m.UpdatedAt)
	return err
}

func (r *EvidenceMappingRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.EvidenceMapping, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, evidence_id, control_id, notes, mapped_by, created_at, updated_at
		 FROM compliance_evidence_mappings WHERE id = $1`, id)
	return scanEvidenceMapping(row)
}

func (r *EvidenceMappingRepo) ListByControl(ctx context.Context, controlID uuid.UUID) ([]*domain.EvidenceMapping, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, evidence_id, control_id, notes, mapped_by, created_at, updated_at
		 FROM compliance_evidence_mappings WHERE control_id = $1`, controlID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEvidenceMappings(rows)
}

func (r *EvidenceMappingRepo) ListByEvidence(ctx context.Context, evidenceID uuid.UUID) ([]*domain.EvidenceMapping, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, evidence_id, control_id, notes, mapped_by, created_at, updated_at
		 FROM compliance_evidence_mappings WHERE evidence_id = $1`, evidenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEvidenceMappings(rows)
}

func (r *EvidenceMappingRepo) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.EvidenceMapping, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, evidence_id, control_id, notes, mapped_by, created_at, updated_at
		 FROM compliance_evidence_mappings WHERE tenant_id = $1`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEvidenceMappings(rows)
}

func (r *EvidenceMappingRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM compliance_evidence_mappings WHERE id = $1`, id)
	return err
}

func (r *EvidenceMappingRepo) GetByEvidenceAndControl(ctx context.Context, evidenceID, controlID uuid.UUID) (*domain.EvidenceMapping, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, evidence_id, control_id, notes, mapped_by, created_at, updated_at
		 FROM compliance_evidence_mappings WHERE evidence_id = $1 AND control_id = $2`, evidenceID, controlID)
	return scanEvidenceMapping(row)
}

// --- QuestionMappingRepo ---

type QuestionMappingRepo struct {
	pool *pgxpool.Pool
}

func NewQuestionMappingRepo(pool *pgxpool.Pool) *QuestionMappingRepo {
	return &QuestionMappingRepo{pool: pool}
}

func (r *QuestionMappingRepo) Create(ctx context.Context, m *domain.QuestionMapping) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO compliance_question_mappings (id, tenant_id, question_id, control_id, mapped_by, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		m.ID, m.TenantID, m.QuestionID, m.ControlID, m.MappedBy, m.CreatedAt)
	return err
}

func (r *QuestionMappingRepo) ListByControl(ctx context.Context, controlID uuid.UUID) ([]*domain.QuestionMapping, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, question_id, control_id, mapped_by, created_at
		 FROM compliance_question_mappings WHERE control_id = $1`, controlID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanQuestionMappings(rows)
}

func (r *QuestionMappingRepo) ListByQuestion(ctx context.Context, questionID uuid.UUID) ([]*domain.QuestionMapping, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, question_id, control_id, mapped_by, created_at
		 FROM compliance_question_mappings WHERE question_id = $1`, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanQuestionMappings(rows)
}

// --- scanners ---

func scanFramework(row pgx.Row) (*domain.Framework, error) {
	var f domain.Framework
	err := row.Scan(&f.ID, &f.Name, &f.Slug, &f.Description, &f.Version, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &f, nil
}

func scanFrameworks(rows pgx.Rows) ([]*domain.Framework, error) {
	var items []*domain.Framework
	for rows.Next() {
		var f domain.Framework
		if err := rows.Scan(&f.ID, &f.Name, &f.Slug, &f.Description, &f.Version, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &f)
	}
	if items == nil {
		return []*domain.Framework{}, nil
	}
	return items, nil
}

func scanControl(row pgx.Row) (*domain.Control, error) {
	var c domain.Control
	err := row.Scan(&c.ID, &c.FrameworkID, &c.ControlID, &c.Name, &c.Description, &c.Category, &c.RiskDescription, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func scanControls(rows pgx.Rows) ([]*domain.Control, error) {
	var items []*domain.Control
	for rows.Next() {
		var c domain.Control
		if err := rows.Scan(&c.ID, &c.FrameworkID, &c.ControlID, &c.Name, &c.Description, &c.Category, &c.RiskDescription, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &c)
	}
	if items == nil {
		return []*domain.Control{}, nil
	}
	return items, nil
}

func scanEvidenceMapping(row pgx.Row) (*domain.EvidenceMapping, error) {
	var m domain.EvidenceMapping
	err := row.Scan(&m.ID, &m.TenantID, &m.EvidenceID, &m.ControlID, &m.Notes, &m.MappedBy, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}

func scanEvidenceMappings(rows pgx.Rows) ([]*domain.EvidenceMapping, error) {
	var items []*domain.EvidenceMapping
	for rows.Next() {
		var m domain.EvidenceMapping
		if err := rows.Scan(&m.ID, &m.TenantID, &m.EvidenceID, &m.ControlID, &m.Notes, &m.MappedBy, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &m)
	}
	if items == nil {
		return []*domain.EvidenceMapping{}, nil
	}
	return items, nil
}

func scanQuestionMapping(row pgx.Row) (*domain.QuestionMapping, error) {
	var m domain.QuestionMapping
	err := row.Scan(&m.ID, &m.TenantID, &m.QuestionID, &m.ControlID, &m.MappedBy, &m.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}

func scanQuestionMappings(rows pgx.Rows) ([]*domain.QuestionMapping, error) {
	var items []*domain.QuestionMapping
	for rows.Next() {
		var m domain.QuestionMapping
		if err := rows.Scan(&m.ID, &m.TenantID, &m.QuestionID, &m.ControlID, &m.MappedBy, &m.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &m)
	}
	if items == nil {
		return []*domain.QuestionMapping{}, nil
	}
	return items, nil
}
