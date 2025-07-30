--
--
--

BEGIN;

-- create the enumeration
CREATE TYPE mtype AS ENUM ('download', 'view');

-- create page_metrics table
CREATE TABLE page_metrics (
   id          serial PRIMARY KEY,
   metric_type mtype,
   namespace   VARCHAR( 32 ) NOT NULL DEFAULT '' ,
   oid         VARCHAR( 64 ) NOT NULL DEFAULT '',
   file_id     VARCHAR( 64 ) NOT NULL DEFAULT '',
   source_ip   VARCHAR( 32 ) NOT NULL DEFAULT '' ,
   referrer    VARCHAR( 255 ) NOT NULL DEFAULT '',
   user_agent  VARCHAR( 32 ) NOT NULL DEFAULT '' ,
   accept_lang VARCHAR( 32 ) NOT NULL DEFAULT '',

   event_time timestamp with time zone NOT NULL,
   created_at timestamp with time zone NOT NULL DEFAULT NOW()
);

-- create the namespace/oid key index
CREATE INDEX page_metrics_key_idx ON page_metrics(namespace, oid);

COMMIT;

--
-- end of file
--
