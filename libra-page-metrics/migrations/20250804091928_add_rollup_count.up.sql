--
--
--

BEGIN;

-- add the rollup_count column
ALTER TABLE page_metrics
ADD COLUMN rollup_count integer NOT NULL DEFAULT 1;

COMMIT;

--
-- end of file
--
