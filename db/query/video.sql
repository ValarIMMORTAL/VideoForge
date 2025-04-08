-- name: InsertVideo :one
insert into videos(
                   title,
                   url,
                   duration,
                   user_id
) values (
          $1,$2,$3,$4
         )  RETURNING *;

-- name: GetVideosByUid :many
select * from videos
where user_id = $1
and delete_at is null;

-- name: GetVideosById :one
select * from videos
where id = $1;