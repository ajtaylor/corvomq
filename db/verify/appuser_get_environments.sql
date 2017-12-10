-- Verify appuser_get_environments

BEGIN;

SELECT 'corvomq.appuser_get_environments(integer)'::regprocedure;

ROLLBACK;
