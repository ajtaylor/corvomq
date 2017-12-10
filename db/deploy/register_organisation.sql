-- Deploy register_organisation
-- requires: appuser
-- requires: organisation

BEGIN;

CREATE OR REPLACE FUNCTION corvomq.register_organisation
  (OUT appuser_id INTEGER
    , OUT organisation_id INTEGER
    , _firstname VARCHAR
    , _lastname VARCHAR
    , _email_address VARCHAR
    , _hashed_password BYTEA
    , _organisation_name VARCHAR
    , _broadcast_key CHAR(36)
    )
AS
$$
DECLARE
  _appuser_id INTEGER;
  _organisation_id INTEGER;
BEGIN

  INSERT INTO corvomq.appuser
    (email_address
      , hashed_password
      , firstname
      , lastname
      )
    VALUES (_email_address
            , _hashed_password
            , _firstname
            , _lastname
            )
    RETURNING id INTO _appuser_id;

  UPDATE corvomq.appuser
  SET created_by = _appuser_id
      , last_updated_by = _appuser_id
      , last_updated_at = NOW()
  WHERE id = _appuser_id;

  INSERT INTO corvomq.organisation
    (name
      , primary_contact_id
      , broadcast_key
      , free_environment_id
      , created_by
      , last_updated_by)
    SELECT _organisation_name
            , _appuser_id
            , _broadcast_key
            , COALESCE(MIN(free_environment_id) + 1, 1)
            , _appuser_id
            , _appuser_id
    FROM corvomq.organisation
    RETURNING id INTO _organisation_id;

  UPDATE corvomq.appuser
  SET organisation_id = _organisation_id
      , last_updated_by = _appuser_id
      , last_updated_at = NOW()
  WHERE id = _appuser_id;

  appuser_id = _appuser_id;
  organisation_id = _organisation_id;

END;
$$
LANGUAGE plpgsql VOLATILE;

COMMIT;
