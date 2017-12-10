-- Revert environment_status

BEGIN;

DROP TABLE corvomq.environment_status;

COMMIT;
