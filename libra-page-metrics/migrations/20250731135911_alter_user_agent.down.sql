--
--
--

BEGIN;

-- restore previous field size
ALTER TABLE page_metrics
ALTER COLUMN user_agent
TYPE VARCHAR(32);

COMMIT;

--
-- end of file
--
