-- Deploy appuser_get_environments
-- requires: appschema

BEGIN;

CREATE OR REPLACE FUNCTION corvomq.appuser_get_environments
  (_appuser_id INTEGER)
  RETURNS TABLE (id INTEGER
                  , name CHARACTER VARYING(200)
                  , server CHARACTER VARYING(20)
                  , infrastructure CHARACTER VARYING(20)
                  , url TEXT
                  , tls_enabled BOOLEAN
                  , created_at TIMESTAMP(3))
AS
$$
BEGIN

  RETURN QUERY
  SELECT
    ef.id
    , 'POC'::CHARACTER VARYING AS "name"
    , 'FREE'::CHARACTER VARYING AS "server"
    , 'FREE'::CHARACTER VARYING AS "infrastructure"
    , 'nats://' || ef.username || ':' || ef.password || '@' || ef.host || ':' || ef.port AS "url"
    , false
    , o.created_at
    FROM corvomq.appuser a
      INNER JOIN corvomq.organisation o
        ON a.organisation_id = o.id
      INNER JOIN corvomq.environment_free_config ef
        ON o.free_environment_id = ef.id
    WHERE a.id = _appuser_id
  UNION
  SELECT
    e.id
    , e.name
    , e.server
    , e.infrastructure
    , e.url
    , e.tls_enabled
    , e.created_at
    FROM corvomq.appuser a
      INNER JOIN corvomq.environment e
        ON a.organisation_id = e.organisation_id
    WHERE a.id = _appuser_id;

END;
$$
LANGUAGE plpgsql VOLATILE;

COMMIT;

