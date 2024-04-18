--
--
--

Begin;

-- create audit table
CREATE TABLE audits (
   id         serial PRIMARY KEY,
   who        VARCHAR( 32 ) NOT NULL DEFAULT '' ,
   oid        VARCHAR( 64 ) NOT NULL DEFAULT '',
   field_name VARCHAR( 32 ) NOT NULL DEFAULT '',
   before     TEXT NOT NULL DEFAULT '',
   after      TEXT NOT NULL DEFAULT '',

   created_at timestamp DEFAULT NOW()
);

-- create the namespace/oid key index
CREATE INDEX audits_key_idx ON audits(who, oid);


COMMIT;

--
--
--
