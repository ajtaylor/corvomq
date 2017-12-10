-- Verify environment

BEGIN;

SELECT
  id
  , organisation_id
  , name
  , server
  , infrastructure
  , url
  , status_id
  , tls_enabled
  , created_at
  , created_by
  , last_updated_at
  , last_updated_by
FROM
  corvomq.environment
WHERE
  FALSE;

ROLLBACK;
