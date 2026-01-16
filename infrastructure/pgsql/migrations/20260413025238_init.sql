-- Create "applications" table
CREATE TABLE "applications" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "redirect_uris" text[] NOT NULL,
  "origin_uris" text[] NOT NULL,
  "client_id" character varying(255) NOT NULL,
  "client_secret" character varying(255) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_applications_client_id" UNIQUE ("client_id")
);
-- Create index "idx_applications_client_id" to table: "applications"
CREATE INDEX "idx_applications_client_id" ON "applications" ("client_id");
-- Create "authentication_providers" table
CREATE TABLE "authentication_providers" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "application_id" uuid NOT NULL,
  "provider" character varying(255) NOT NULL,
  "client_id" character varying(255) NOT NULL,
  "client_secret" character varying(255) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_authentication_providers_application" FOREIGN KEY ("application_id") REFERENCES "applications" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_authentication_providers_provider" to table: "authentication_providers"
CREATE INDEX "idx_authentication_providers_provider" ON "authentication_providers" ("provider");
-- Create index "idx_client_application" to table: "authentication_providers"
CREATE UNIQUE INDEX "idx_client_application" ON "authentication_providers" ("application_id", "client_id");
-- Create "authentication_requests" table
CREATE TABLE "authentication_requests" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "application_id" uuid NOT NULL,
  "provider_id" uuid NULL,
  "response_type" character varying(255) NOT NULL,
  "redirect_uri" text NOT NULL,
  "state" character varying(255) NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "expires_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP + '00:30:00'::interval),
  "completed_at" timestamp NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_authentication_requests_application" FOREIGN KEY ("application_id") REFERENCES "applications" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "application_id" uuid NOT NULL,
  "provider_id" uuid NOT NULL,
  "provider_user_id" character varying(255) NOT NULL,
  "email" character varying(255) NOT NULL,
  "name" character varying(255) NULL,
  "avatar_url" text NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_application" FOREIGN KEY ("application_id") REFERENCES "applications" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_users_provider" FOREIGN KEY ("provider_id") REFERENCES "authentication_providers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_app_provider_user" to table: "users"
CREATE UNIQUE INDEX "idx_app_provider_user" ON "users" ("application_id", "provider_user_id");
-- Create index "idx_users_email" to table: "users"
CREATE INDEX "idx_users_email" ON "users" ("email");
-- Create index "idx_users_provider_user_id" to table: "users"
CREATE INDEX "idx_users_provider_user_id" ON "users" ("provider_user_id");
-- Create "authorization_codes" table
CREATE TABLE "authorization_codes" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "code" character varying(255) NOT NULL,
  "authentication_request_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "used_at" timestamp NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "expires_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP + '00:10:00'::interval),
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_authorization_codes_code" UNIQUE ("code"),
  CONSTRAINT "fk_authorization_codes_authentication_request" FOREIGN KEY ("authentication_request_id") REFERENCES "authentication_requests" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_authorization_codes_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_authorization_codes_code" to table: "authorization_codes"
CREATE INDEX "idx_authorization_codes_code" ON "authorization_codes" ("code");
