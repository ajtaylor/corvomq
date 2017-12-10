-- Verify environment_status

BEGIN;

SELECT
  id
  , code
  , display
  , created_at
  , created_by
  , last_updated_at
  , last_updated_by
FROM
  corvomq.environment_status
WHERE
  FALSE;

ROLLBACK;
