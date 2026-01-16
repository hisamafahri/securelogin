-- Create "sessions" table
CREATE TABLE "sessions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "token" character varying(255) NOT NULL,
  "authentication_request_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "expires_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP + '01:00:00'::interval),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_sessions_token" UNIQUE ("token"),
  CONSTRAINT "fk_sessions_authentication_request" FOREIGN KEY ("authentication_request_id") REFERENCES "authentication_requests" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_sessions_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_sessions_token" to table: "sessions"
CREATE INDEX "idx_sessions_token" ON "sessions" ("token");
