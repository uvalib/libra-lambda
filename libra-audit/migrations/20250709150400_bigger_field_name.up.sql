BEGIN;

-- Increase the size of the field_name column
ALTER TABLE audits
ALTER COLUMN field_name
TYPE VARCHAR( 127 );

COMMIT;