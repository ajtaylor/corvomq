-- Deploy login_appuser
-- requires: appuser
-- requires: organisation

BEGIN;

CREATE OR REPLACE FUNCTION corvomq.login_appuser
  (_email_address VARCHAR)
  RETURNS TABLE (appuser_id INTEGER
                  , hashed_password BYTEA)
AS
$$
BEGIN

  RETURN QUERY
  SELECT
    u.id
    , u.hashed_password
    FROM corvomq.appuser u
    WHERE u.email_address = _email_address;

END;
$$
LANGUAGE plpgsql VOLATILE;

END;
