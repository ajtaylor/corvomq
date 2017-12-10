-- Revert appuser_get_environments

BEGIN;

DROP FUNCTION corvomq.appuser_get_environments (INTEGER);

COMMIT;
