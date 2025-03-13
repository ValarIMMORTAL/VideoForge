-- name: CreateCopy :one
INSERT INTO copywriting(
    title,
    source,
    content,
    date
) values (
    $1,$2,$3,$4
) returning *;

-- name: GetCopy :one
select * from copywriting
where id = $1
  and delete_at = null
    limit 1;

-- name: ListCopies :many
select * from copywriting
where date = $1
and deleta_at = null
order by id
limit $2
offset $3;

-- name: DeleteCopy :exec
update copywriting
set delete_at = $1
where id = $2;