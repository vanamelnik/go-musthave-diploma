CREATE TYPE "order_status" AS ENUM (
  'NEW',
  'REGISTERED',
  'INVALID',
  'PROCESSING',
  'PROCESSED'
);

CREATE TABLE "users" (
  "id" uuid UNIQUE PRIMARY KEY,
  "login" text UNIQUE NOT NULL,
  "password_hash" text NOT NULL,
  "gpoints_balance" decimal default 0,
  "remember_token" VARCHAR(50) NOT NULL DEFAULT '',
  "created_at" timestamp
);

CREATE TABLE "orders" (
  "id" text UNIQUE PRIMARY KEY NOT NULL,
  "user_id" uuid,
  "status" order_status NOT NULL,
  "accrual_points" decimal,
  "uploaded_at" timestamp
);

CREATE TABLE "accruals_log" (
  "order_id" text UNIQUE NOT NULL,
  "user_id" uuid,
  "processed" boolean DEFAULT false,
  "sum" decimal
);

CREATE TABLE "withdrawals_log" (
  "order_id" text UNIQUE NOT NULL,
  "user_id" uuid,
  "sum" decimal,
  "status" order_status, 
  "processed_at" timestamp
);

--ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

--ALTER TABLE "accruals_log" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

--ALTER TABLE "withdrawals_log" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
