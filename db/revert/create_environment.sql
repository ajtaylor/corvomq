-- Revert create_environment

BEGIN;

DROP FUNCTION corvomq.create_environment
  (INTEGER, VARCHAR, VARCHAR, VARCHAR, VARCHAR, BOOLEAN);

COMMIT;
