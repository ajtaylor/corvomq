-- Deploy populate_environment_status
-- requires: environment_status

BEGIN;

INSERT INTO corvomq.environment_status
    (id,
      code,
      display)
VALUES (1, 'BUILDING', 'Building...'),
        (2, 'RUNNING', 'Running'),
        (3, 'STOPPED', 'Stopped');

COMMIT;
