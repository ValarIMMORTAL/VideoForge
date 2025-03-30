CREATE TABLE "oauth2_tokens" (
                                 "id" bigint PRIMARY KEY,
                                 "user_id" integer NOT NULL,
                                 "provider" varchar NOT NULL,
                                 "api" varchar NOT NULL,
                                 "refresh_token" text NOT NULL,
                                 "created_at" timestamptz NOT NULL DEFAULT (now()),
                                 "delete_at" timestamp default null
);

ALTER TABLE "oauth2_tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");