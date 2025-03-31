CREATE TABLE "oauth2_tokens" (
                                 "id" BIGSERIAL PRIMARY KEY,
                                 "user_id" BIGINT NOT NULL,
                                 "provider" varchar NOT NULL,
                                 "api" varchar NOT NULL,
                                 "access_token" text NOT NULL,
                                 "token_type" varchar not null,
                                 "refresh_token" text NOT NULL,
                                 "expiry" TIMESTAMPTZ not null ,
                                 "created_at" timestamptz NOT NULL DEFAULT (now()),
                                 "delete_at" timestamp default null
);

ALTER TABLE "oauth2_tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");