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

-- name: CreateMultipleCopy :exec
INSERT INTO copywriting(
    title,
    source,
    content,
    date
)
select
    unnest($1::text[]),
    unnest($2::text[]),
    unnest($3::text[]),
    unnest($4::timestamp[]);