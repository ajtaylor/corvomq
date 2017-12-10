-- Deploy environment_free_config
-- requires: appschema

BEGIN;

CREATE TABLE IF NOT EXISTS corvomq.environment_free_config
(
  id SERIAL NOT NULL
  , host CHARACTER VARYING(200) NOT NULL
  , port INTEGER NOT NULL
  , username CHARACTER VARYING(15) NOT NULL
  , password CHARACTER VARYING(15) NOT NULL
  , message_prefix CHARACTER VARYING(15) NOT NULL
  , CONSTRAINT environment_free_config_pk PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

COMMIT;
