-- Deploy create_environment
-- requires: environment
-- requires: environment_status

BEGIN;

CREATE OR REPLACE FUNCTION corvomq.create_environment
  (OUT environment_id INTEGER
    , _appuser_id INTEGER
    , _name VARCHAR
    , _server VARCHAR
    , _infrastructure VARCHAR
    , _urlsegment VARCHAR
    , _tls_enabled BOOLEAN
    )
AS
$$
DECLARE
  _environment_id INTEGER;
BEGIN

  INSERT INTO corvomq.environment
    (organisation_id
      , name
      , server
      , infrastructure
      , url
      , tls_enabled
      , status_id
      , created_by
      , last_updated_by
      )
    SELECT organisation_id
            , _name
            , _server
            , _infrastructure
            , 'http://' || _urlsegment || '.corvomq.com'
            , _tls_enabled
            , es.id
            , a.id
            , a.id
    FROM corvomq.appuser a
      CROSS JOIN corvomq.environment_status es
    WHERE a.id = _appuser_id
      AND es.code = 'BUILDING'
    RETURNING id INTO _environment_id;

  environment_id = _environment_id;

END;
$$
LANGUAGE plpgsql VOLATILE;

COMMIT;
