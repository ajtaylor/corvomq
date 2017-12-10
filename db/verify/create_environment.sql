-- Verify create_environment

BEGIN;

SELECT 'corvomq.create_environment(integer, varchar, varchar, varchar, varchar, boolean)'::regprocedure;

ROLLBACK;
