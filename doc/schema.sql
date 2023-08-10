-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-08-03T08:18:56.743Z

CREATE TABLE "user" (
  "name" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "account" (
  "id" INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "account_id" varchar NOT NULL,
  "username" varchar NOT NULL,
  "created" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "snippets" (
  "id" INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "user_id" integer NOT NULL,
  "title" varchar,
  "content" varchar NOT NULL,
  "expires" timestamptz NOT NULL,
  "created" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "session" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created" timestamptz NOT NULL DEFAULT 'now()'
);

ALTER TABLE "account" ADD FOREIGN KEY ("account_id") REFERENCES "user" ("name");

ALTER TABLE "snippets" ADD FOREIGN KEY ("user_id") REFERENCES "account" ("id");

ALTER TABLE "session" ADD FOREIGN KEY ("name") REFERENCES "user" ("name");