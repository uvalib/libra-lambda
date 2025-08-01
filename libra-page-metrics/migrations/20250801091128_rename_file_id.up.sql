--
--
--

BEGIN;

-- rename the column
ALTER TABLE page_metrics
RENAME COLUMN file_id TO target_id;

COMMIT;

--
-- end of file
--
