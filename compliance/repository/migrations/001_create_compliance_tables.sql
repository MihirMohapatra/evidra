CREATE TABLE compliance_frameworks (
    id          UUID         PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT         NOT NULL DEFAULT '',
    version     VARCHAR(50)  NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE compliance_controls (
    id               UUID         PRIMARY KEY,
    framework_id     UUID         NOT NULL REFERENCES compliance_frameworks(id) ON DELETE CASCADE,
    control_id       VARCHAR(50)  NOT NULL,
    name             VARCHAR(255) NOT NULL,
    description      TEXT         NOT NULL DEFAULT '',
    category         VARCHAR(30)  NOT NULL,
    risk_description TEXT         NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE(framework_id, control_id)
);

CREATE TABLE compliance_evidence_mappings (
    id          UUID         PRIMARY KEY,
    tenant_id   UUID         NOT NULL,
    evidence_id UUID         NOT NULL,
    control_id  UUID         NOT NULL REFERENCES compliance_controls(id) ON DELETE CASCADE,
    notes       TEXT         NOT NULL DEFAULT '',
    mapped_by   UUID         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE(evidence_id, control_id)
);

CREATE TABLE compliance_question_mappings (
    id          UUID         PRIMARY KEY,
    tenant_id   UUID         NOT NULL,
    question_id UUID         NOT NULL,
    control_id  UUID         NOT NULL REFERENCES compliance_controls(id) ON DELETE CASCADE,
    mapped_by   UUID         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_controls_framework ON compliance_controls(framework_id);
CREATE INDEX idx_ev_mappings_control ON compliance_evidence_mappings(control_id);
CREATE INDEX idx_ev_mappings_evidence ON compliance_evidence_mappings(evidence_id);
CREATE INDEX idx_ev_mappings_tenant ON compliance_evidence_mappings(tenant_id);
CREATE INDEX idx_q_mappings_control ON compliance_question_mappings(control_id);
CREATE INDEX idx_q_mappings_question ON compliance_question_mappings(question_id);

-- Seed common compliance frameworks
INSERT INTO compliance_frameworks (id, name, slug, description, version) VALUES
    (gen_random_uuid(), 'SOC 2', 'soc2', 'Service Organization Control 2 - Trust Services Criteria', '2023'),
    (gen_random_uuid(), 'ISO 27001', 'iso27001', 'Information Security Management System standard', '2022'),
    (gen_random_uuid(), 'NIST 800-53', 'nist80053', 'Security and Privacy Controls for Information Systems', 'Rev 5'),
    (gen_random_uuid(), 'PCI DSS', 'pci-dss', 'Payment Card Industry Data Security Standard', '4.0'),
    (gen_random_uuid(), 'HIPAA', 'hipaa', 'Health Insurance Portability and Accountability Act', '2023'),
    (gen_random_uuid(), 'FedRAMP', 'fedramp', 'Federal Risk and Authorization Management Program', 'Rev 5')
ON CONFLICT (slug) DO NOTHING;
