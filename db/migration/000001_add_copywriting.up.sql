CREATE TABLE "copywriting" (
                               "id" BIGSERIAL PRIMARY KEY,
                               "source" varchar  not null,
                               "title" varchar not null ,
                               "content" varchar not null ,
                               "date" timestamp not null ,
                               "created_at" timestamptz NOT NULL DEFAULT (now()),
                               "delete_at" timestamp default null
);
