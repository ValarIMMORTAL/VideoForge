-- name: InsertThirdKey :one
insert into third_keys(
    name,
    ak,
    sk
) values (
          $1,$2,$3
         ) RETURNING *;


-- name: GetThirdKeyByName :one
select *
from third_keys
where name = $1
and delete_at is null;
