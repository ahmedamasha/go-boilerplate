DROP TABLE IF EXISTS withdrawal_requests;

ALTER TABLE events
    DROP COLUMN IF EXISTS amount_withdrawn,
    DROP COLUMN IF EXISTS status;

DROP INDEX IF EXISTS idx_events_status;
