-- Deploy organisation
-- requires: appschema

BEGIN;

CREATE TABLE IF NOT EXISTS corvomq.organisation
(
  id serial NOT NULL
  , name CHARACTER VARYING(200) NOT NULL
  , primary_contact_id INTEGER
  , broadcast_key CHAR(36)
  , free_environment_id INTEGER
  , created_at TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT now()
  , created_by INTEGER
  , last_updated_at TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT now()
  , last_updated_by INTEGER
  , CONSTRAINT organisation_pk PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

COMMIT;
