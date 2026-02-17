CREATE TYPE template_state AS ENUM ('DRAFT', 'SUBMITTED', 'REJECTED', 'REVIEWED', 'APPROVED');

CREATE TABLE IF NOT EXISTS contract_templates (
    did VARCHAR(255),
    document_number INT DEFAULT 1,
    version INT DEFAULT 1,

    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    state template_state NOT NULL,

    name VARCHAR(255) NOT NULL,
    description TEXT,

    template_data JSONB DEFAULT '{}'::jsonb,

    CONSTRAINT pk_contract_templates PRIMARY KEY (did, document_number, version),
    CONSTRAINT chk_did_not_empty CHECK (did <> ''),
    CONSTRAINT chk_document_number_positive CHECK (document_number > 0),
    CONSTRAINT chk_version_positive CHECK (version > 0)
);

-- Trigger for updating updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER contract_templates_update_updated_at
    BEFORE UPDATE ON contract_templates
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

------------------------------------------------------------------------------------------------------------------------

CREATE TYPE review_task_state AS ENUM ('OPEN', 'APPROVED', 'REJECTED');

CREATE TABLE IF NOT EXISTS contract_templates_review_task
(
    id              BIGSERIAL PRIMARY KEY,

    did             VARCHAR(255) CHECK (did <> ''),
    document_number INT NOT NULL,
    version         INT NOT NULL,

    state review_task_state NOT NULL,
    reviewer VARCHAR(255) CHECK (reviewer <> '' AND reviewer IS NOT NULL),

    created_by      VARCHAR(255) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_review_task_contract_template
        FOREIGN KEY (did, document_number, version)
        REFERENCES contract_templates(did, document_number, version)
);

------------------------------------------------------------------------------------------------------------------------

CREATE TYPE approval_task_state AS ENUM ('OPEN', 'APPROVED', 'REJECTED', 'RESUBMITTED');

CREATE TABLE IF NOT EXISTS contract_templates_approval_task
(
    id              BIGSERIAL PRIMARY KEY,

    did             VARCHAR(255) CHECK (did <> ''),
    document_number INT NOT NULL,
    version         INT NOT NULL,

    state approval_task_state NOT NULL,
    approver VARCHAR(255) CHECK (approver <> '' AND approver IS NOT NULL),

    created_by      VARCHAR(255) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_review_task_contract_template
        FOREIGN KEY (did, document_number, version)
            REFERENCES contract_templates(did, document_number, version)
);