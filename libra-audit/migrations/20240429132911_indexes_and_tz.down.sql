BEGIN;


-- Revert to no tz
ALTER TABLE audits
ALTER COLUMN created_at
TYPE TIMESTAMP;


DROP INDEX IF EXISTS audits_who_idx;

COMMIT;
