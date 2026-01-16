-- Modify "authentication_requests" table
ALTER TABLE "authentication_requests" ADD COLUMN "scopes" text[] NOT NULL;
