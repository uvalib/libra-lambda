BEGIN;

-- Add tz to timestamp
ALTER TABLE audits
ALTER COLUMN created_at
TYPE TIMESTAMP WITH TIME ZONE;

SHOW timezone;


-- fix the namespace/oid key index
DROP INDEX IF EXISTS audits_key_idx;
CREATE INDEX audits_key_idx ON audits(namespace, oid);

-- add who index
CREATE INDEX audits_who_idx on audits who;


COMMIT;