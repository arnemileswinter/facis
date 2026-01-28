CREATE TABLE IF NOT EXISTS contract_templates (
    did VARCHAR(255) PRIMARY KEY CHECK (did <> '' AND did IS NOT NULL),
    document_number INT NOT NULL,
    version INT NOT NULL,
    state VARCHAR(16) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    meta_data JSONB
);
