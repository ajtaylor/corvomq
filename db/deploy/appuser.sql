-- Deploy appuser
-- requires: appschema organisation

BEGIN;

CREATE TABLE IF NOT EXISTS corvomq.appuser
(
  id serial NOT NULL
  , email_address CHARACTER VARYING(200) NOT NULL
  , hashed_password BYTEA NOT NULL
  , firstname CHARACTER VARYING(200)
  , lastname CHARACTER VARYING(200)
  , organisation_id INTEGER
  , created_at TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT now()
  , created_by INTEGER
  , last_updated_at TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT now()
  , last_updated_by INTEGER
  , CONSTRAINT appuser_pk PRIMARY KEY (id)
  , CONSTRAINT appuser_uk_email_address UNIQUE (email_address)
  , CONSTRAINT fk_appuser_organisation
      FOREIGN KEY (organisation_id) REFERENCES corvomq.organisation(id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
)
WITH (
  OIDS=FALSE
);

COMMIT;
