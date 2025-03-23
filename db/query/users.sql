-- name: CreateUser :one
insert into users (username,
                   hashed_password,
                   email) values (
    $1,$2,$3
    ) RETURNING *;

-- name: GetUser :one
select * from users
where id = $1 limit 1;

-- name: GetUserByName :one
select * from users
where username = $1 limit 1;

-- name: UpdateUser :one
UPDATE users
SET
    username = COALESCE(sqlc.narg(username), username),
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
    email = COALESCE(sqlc.narg(email), email)
where
    username = sqlc.arg(username)
    RETURNING *;
