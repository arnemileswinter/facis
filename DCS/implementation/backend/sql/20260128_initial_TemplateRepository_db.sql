CREATE TYPE template_state AS ENUM ('DRAFT', 'SUBMITTED', 'REVIEWED', 'APPROVED');

CREATE TABLE IF NOT EXISTS contract_templates (
                                                  did VARCHAR(255) PRIMARY KEY CHECK (did <> '' AND did IS NOT NULL),

                                                  created_by VARCHAR(255) NOT NULL,
                                                  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                                                  updated_by VARCHAR(255),
                                                  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                                  document_number INT DEFAULT 1 CHECK (document_number > 0),
                                                  version INT DEFAULT 1 CHECK (version > 0),

                                                  state template_state NOT NULL,

                                                  name VARCHAR(255) NOT NULL,
                                                  description TEXT,
                                                  meta_data JSONB DEFAULT '{}'::jsonb
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

-- Indizes
CREATE INDEX idx_contract_templates_state ON contract_templates(state);
CREATE INDEX idx_contract_templates_created_by ON contract_templates(created_by);
CREATE INDEX idx_contract_templates_updated_by ON contract_templates(updated_by);

------------------------------------------------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS contract_templates_review_task
(
    did             VARCHAR(255) PRIMARY KEY CHECK (did <> '' AND did IS NOT NULL),

    created_by      VARCHAR(255)   NOT NULL,
    created_at      TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_by      VARCHAR(255),
    updated_at      TIMESTAMP               DEFAULT CURRENT_TIMESTAMP,

    document_number INT                     DEFAULT 1 CHECK (document_number > 0),
    version         INT                     DEFAULT 1 CHECK (version > 0),

    state           template_state NOT NULL
);

CREATE TRIGGER contract_templates_review_task_update_updated_at
    BEFORE UPDATE ON contract_templates_review_task
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();