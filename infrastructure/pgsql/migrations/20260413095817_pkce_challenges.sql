-- Modify "authentication_requests" table
ALTER TABLE "authentication_requests" ADD COLUMN "code_challenge" text NULL, ADD COLUMN "code_challenge_method" character varying(255) NULL;
