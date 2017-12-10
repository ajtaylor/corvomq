-- Revert populate_environment_status

BEGIN;

TRUNCATE TABLE corvomq.environment_status
  RESTART IDENTITY;

COMMIT;
