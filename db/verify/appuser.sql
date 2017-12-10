-- Verify appusers

BEGIN;

SELECT
  id
  , email_address
  , hashed_password
  , firstname
  , lastname
  , organisation_id
  , created_at
  , created_by
  , last_updated_at
  , last_updated_by
FROM
  corvomq.appuser
WHERE
  FALSE;

ROLLBACK;
