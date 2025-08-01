--
--
--

BEGIN;

-- rename the column
ALTER TABLE page_metrics
RENAME COLUMN target_id TO file_id;

COMMIT;

--
-- end of file
--
