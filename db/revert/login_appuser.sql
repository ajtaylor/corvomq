-- Revert login_appuser

BEGIN;

DROP FUNCTION corvomq.login_appuser (VARCHAR);

COMMIT;
