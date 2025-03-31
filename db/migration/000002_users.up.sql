CREATE TABLE "users" (
                         "id" BIGSERIAL PRIMARY KEY,
                         "username" varchar not null ,
                         "hashed_password" varchar not null ,
                         "email" varchar not null ,
                         "created_at" timestamptz NOT NULL DEFAULT (now())
);