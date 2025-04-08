CREATE TABLE "videos" (
                                 "id" BIGSERIAL PRIMARY KEY,
                                 "title" varchar not null ,
                                 "url" varchar not null ,
                                 "duration" integer not null ,
                                 "user_id" BIGINT NOT NULL,
                                 "subscribe" BIGINT not null ,
                                 "created_at" timestamptz NOT NULL DEFAULT (now()),
                                 "delete_at" timestamp default null
);

ALTER TABLE "videos" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");