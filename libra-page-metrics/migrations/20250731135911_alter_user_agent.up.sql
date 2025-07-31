--
--
--

BEGIN;

-- increase field size
ALTER TABLE page_metrics
ALTER COLUMN user_agent
TYPE VARCHAR(255);

COMMIT;

--
-- end of file
--
