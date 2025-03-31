-- name: InsertOauth2Token :one
insert into oauth2_tokens(
                          user_id,
                          provider,
                          api,
                          access_token,
                          token_type,
                          refresh_token,
                          expiry
) values (
          $1,$2,$3,$4,$5,$6,$7
         ) RETURNING *;

-- name: GetOauth2Token :one
select id, user_id, provider,api,access_token,token_type,refresh_token,expiry
from oauth2_tokens
where user_id = $1
and provider = $2
and api = $3
and delete_at is null;

-- name: UpdateAccessToken :one
update oauth2_tokens
set access_token = $1, token_type = $2, expiry = $3
where user_id = $4 and provider = $5 and refresh_token = $6 and token_type = $7
RETURNING *;

