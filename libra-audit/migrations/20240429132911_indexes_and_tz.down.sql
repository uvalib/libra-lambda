BEGIN;


-- Revert to no tz
ALTER TABLE audits
ALTER COLUMN created_at
TYPE TIMESTAMP;


DROP INDEX audits_who_idx;

COMMIT;
