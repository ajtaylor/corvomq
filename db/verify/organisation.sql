-- Verify organisation

BEGIN;

SELECT
  id
  , name
  , primary_contact_id
  , broadcast_key
  , free_environment_id
  , created_at
  , created_by
  , last_updated_at
  , last_updated_by
FROM
  corvomq.organisation
WHERE
  FALSE;

ROLLBACK;
