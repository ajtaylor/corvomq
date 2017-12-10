-- Verify register_organisation

BEGIN;

SELECT 'corvomq.register_organisation(varchar, varchar, varchar, bytea, varchar, char)'::regprocedure;

ROLLBACK;
