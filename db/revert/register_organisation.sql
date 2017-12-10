-- Revert register_organisation

BEGIN;

DROP FUNCTION corvomq.register_organisation
  (VARCHAR, VARCHAR, VARCHAR, BYTEA, VARCHAR, CHAR(36));

COMMIT;
