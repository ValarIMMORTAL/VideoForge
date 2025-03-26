CREATE TABLE "platforms" (
                             "id" BIGSERIAL PRIMARY KEY,
                             "platform" varchar NOT NULL,
                             "detail" varchar NOT NULL,
                             "created_at" timestamptz NOT NULL DEFAULT (now()),
                             "delete_at" timestamp default null
);