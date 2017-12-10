-- Revert environment_free_config

BEGIN;

DROP TABLE corvomq.environment_free_config;

COMMIT;
