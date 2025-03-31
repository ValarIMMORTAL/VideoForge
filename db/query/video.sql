-- name: InsertVideo :one
insert into videos(
                   title,
                   url,
                   duration,
                   user_id
) values (
          $1,$2,$3,$4
         )  RETURNING *;