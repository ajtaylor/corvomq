-- Revert appusers

BEGIN;

DROP TABLE corvomq.appuser;

COMMIT;
