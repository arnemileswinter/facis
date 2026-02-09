-- Migration: Create Outbox Events Table for CQRS/Event Sourcing
-- Purpose: Persistent storage of events before publishing to NATS
-- This ensures events are never lost even if NATS is unavailable

CREATE TABLE IF NOT EXISTS outbox_events (
    id BIGSERIAL PRIMARY KEY,

    event_type VARCHAR(64) NOT NULL,
    event_data JSONB NOT NULL,

    did VARCHAR(255),

    processed BOOLEAN DEFAULT FALSE,
    processed_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ✅ CRITICAL: Index for Poller query
-- The Poller does: SELECT * FROM outbox_events WHERE processed = FALSE
-- This index must be fast!
CREATE INDEX idx_outbox_events_processed 
    ON outbox_events(processed, created_at)
    WHERE processed = FALSE;

-- Optional: Index for cleanup/audit queries
CREATE INDEX idx_outbox_events_created_at 
    ON outbox_events(created_at);

-- Optional: Index for specific event queries
CREATE INDEX idx_outbox_events_did
    ON outbox_events(did);

-- Optional: Index for event type filtering
CREATE INDEX idx_outbox_events_type 
    ON outbox_events(event_type);

-- ✅ Stored Procedure: Add Event to Outbox (for consistency)
-- Usage in Handler: SELECT add_outbox_event('EventType', json_data, 'did');
CREATE OR REPLACE FUNCTION add_outbox_event(
    p_event_type VARCHAR,
    p_event_data JSONB,
    p_did VARCHAR
) RETURNS BIGINT AS $$
DECLARE
    v_event_id BIGINT;
BEGIN
    INSERT INTO outbox_events (event_type, event_data, did, processed)
    VALUES (p_event_type, p_event_data, p_did, FALSE)
    RETURNING id INTO v_event_id;
    
    RETURN v_event_id;
END;
$$ LANGUAGE plpgsql;

-- ✅ Stored Procedure: Mark Event as Processed
-- Usage: SELECT mark_event_processed(123);
CREATE OR REPLACE FUNCTION mark_event_processed(p_event_id BIGINT)
RETURNS BOOLEAN AS $$
BEGIN
    UPDATE outbox_events 
    SET processed = TRUE, processed_at = CURRENT_TIMESTAMP
    WHERE id = p_event_id;
    
    RETURN FOUND;
END;
$$ LANGUAGE plpgsql;

-- ✅ Stored Procedure: Cleanup old processed events
-- Usage: SELECT cleanup_old_events(interval '7 days');
-- Run periodically (e.g., daily) to prevent table from growing indefinitely
CREATE OR REPLACE FUNCTION cleanup_old_events(p_retention_period INTERVAL DEFAULT '7 days')
RETURNS INTEGER AS $$
DECLARE
    v_deleted_count INTEGER;
BEGIN
    DELETE FROM outbox_events
    WHERE processed = TRUE 
      AND processed_at < CURRENT_TIMESTAMP - p_retention_period;
    
    GET DIAGNOSTICS v_deleted_count = ROW_COUNT;
    RETURN v_deleted_count;
END;
$$ LANGUAGE plpgsql;

-- ✅ View: Unprocessed Events (for monitoring)
CREATE OR REPLACE VIEW v_unprocessed_events AS
SELECT 
    id,
    event_type,
    event_data,
    did,
    created_at,
    CURRENT_TIMESTAMP - created_at AS age
FROM outbox_events
WHERE processed = FALSE
ORDER BY created_at ASC;

-- ✅ View: Event Processing Statistics (for monitoring)
CREATE OR REPLACE VIEW v_outbox_stats AS
SELECT 
    event_type,
    COUNT(*) FILTER (WHERE processed = FALSE) AS unprocessed_count,
    COUNT(*) FILTER (WHERE processed = TRUE) AS processed_count,
    COUNT(*) AS total_count,
    MAX(created_at) AS latest_event
FROM outbox_events
GROUP BY event_type
ORDER BY unprocessed_count DESC;
