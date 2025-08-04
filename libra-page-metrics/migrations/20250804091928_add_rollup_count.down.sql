--
--
--

BEGIN;

-- drop the rollup_count column
ALTER TABLE page_metrics
DROP COLUMN rollup_count;

COMMIT;

--
-- end of file
--
