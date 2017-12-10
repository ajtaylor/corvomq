-- Verify environment_free_config

BEGIN;

SELECT
  id
  , host
  , port
  , username
  , password
  , message_prefix
FROM
  corvomq.environment_free_config
WHERE
  FALSE;

ROLLBACK;
