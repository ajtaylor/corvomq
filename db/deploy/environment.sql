-- Deploy environment
-- requires: appschema

BEGIN;

CREATE TABLE IF NOT EXISTS corvomq.environment
(
  id SERIAL NOT NULL
  , organisation_id INTEGER NOT NULL
  , name CHARACTER VARYING(200) NOT NULL
  , server CHARACTER VARYING(20) NOT NULL
  , infrastructure CHARACTER VARYING(20) NOT NULL
  , url CHARACTER VARYING(200) NOT NULL
  , tls_enabled BOOLEAN NOT NULL
  , status_id INTEGER NOT NULL
  , created_at TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT now()
  , created_by INTEGER
  , last_updated_at TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT now()
  , last_updated_by INTEGER
  , CONSTRAINT environment_pk PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

COMMIT;
