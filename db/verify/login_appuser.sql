-- Verify login_appuser

BEGIN;

SELECT 'corvomq.login_appuser(varchar)'::regprocedure;

ROLLBACK;
