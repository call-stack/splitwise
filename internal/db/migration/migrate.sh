#!/bin/bash
set -e
export PGPASSWORD=$POSTGRES_PASSWORD;
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL

  DROP TABLE IF EXISTS summary ;
  DROP TABLE IF EXISTS user_expense ;
  DROP TABLE IF EXISTS expense ;
  DROP TABLE IF EXISTS user_group ;
  DROP TABLE IF EXISTS groups ;
  DROP TABLE IF EXISTS users ;


  CREATE TABLE "users"
  (
      "uid"          SERIAL PRIMARY KEY,
      "name"         varchar     NOT NULL,
      "email"        varchar     NOT NULL UNIQUE,
      "phone_number" varchar     NOT NULL UNIQUE,
      "created_at"   timestamptz NOT NULL DEFAULT (now())
  );

  CREATE TABLE "groups"
  (
      "gid"        SERIAL PRIMARY KEY,
      "name"       varchar     NOT NULL,
      "created_at" timestamptz NOT NULL DEFAULT (now())
  );

  CREATE TABLE "user_group"
  (
      "uid" integer,
      "gid" integer,
      PRIMARY KEY ("uid", "gid")
  );

  CREATE TABLE "expense"
  (
      "eid"                   SERIAL PRIMARY KEY,
      "uid"                   integer,
      "category"              varchar,
      "is_expense_settled"    bool,
      "amount"                int,
      "people_involved"       int,
      "created_at"            timestamptz NOT NULL DEFAULT (now())
  );

  CREATE TABLE "user_expense"
  (
      "paid_by"            integer,
      "eid"                integer,
      "paid_to"            integer,
      "amount"             float,
      PRIMARY KEY ("paid_by", "eid")
  );
  CREATE TABLE "summary"
  (
      "uid"    integer PRIMARY KEY,
      "amount" float
  );

  ALTER TABLE "user_group"
      ADD CONSTRAINT fk_user_user_group FOREIGN KEY ("uid") REFERENCES users ("uid");
  ALTER TABLE "user_group"
      ADD CONSTRAINT fk_group_user_group FOREIGN KEY ("gid") REFERENCES groups ("gid");
  ALTER TABLE "expense"
      ADD CONSTRAINT fk_user_expense FOREIGN KEY ("uid") REFERENCES users ("uid");
  ALTER TABLE "user_expense"
      ADD CONSTRAINT fk_user_expm_amount_1 FOREIGN KEY ("paid_by") REFERENCES users ("uid");
  ALTER TABLE "user_expense"
      ADD CONSTRAINT fk_user_expm_amount_2 FOREIGN KEY ("paid_to") REFERENCES users ("uid");
  ALTER TABLE "user_expense"
      ADD CONSTRAINT fk_txn_expense_amount_1 FOREIGN KEY ("eid") REFERENCES expense ("eid");
  ALTER TABLE "summary"
      ADD CONSTRAINT fk_user_summary FOREIGN KEY ("uid") REFERENCES users ("uid");
  COMMIT;
EOSQL