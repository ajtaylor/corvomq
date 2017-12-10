-- Deploy environment_status
-- requires: appschema

BEGIN;

CREATE TABLE IF NOT EXISTS corvomq.environment_status
(
  id INTEGER NOT NULL
  , code CHARACTER VARYING(50) NOT NULL
  , display CHARACTER VARYING(200) NOT NULL
  , created_at TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT now()
  , created_by INTEGER
  , last_updated_at TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT now()
  , last_updated_by INTEGER
  , CONSTRAINT environment_status_pk PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

COMMIT;
