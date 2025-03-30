-- name: InsertOauth2Token :one
insert into oauth2_tokens(
                          user_id,
                          provider,
                          api,
                          refresh_token
) values (
          $1,$2,$3,$4
         ) RETURNING *;

-- name: GetOauth2Token :one
select id, user_id, provider,api,refresh_token
from oauth2_tokens
where user_id = $1
and provider = $2
and api = $3
and delete_at is null;