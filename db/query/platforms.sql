-- name: GetPlatforms :many
select id,platform,detail
from platforms
where delete_at is null
order by id;


-- name: GetPlatformsByName :one
select id,platform,detail
from platforms
where delete_at is null
and platform=$1
limit 1;