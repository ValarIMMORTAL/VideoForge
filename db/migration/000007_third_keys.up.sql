CREATE TABLE "third_keys" (
                          "id" BIGSERIAL PRIMARY KEY,
                          "name" varchar not null ,
                          "ak" varchar not null ,
                          "sk" varchar not null ,
                          "created_at" timestamptz NOT NULL DEFAULT (now()),
                          "delete_at" timestamp default null
);