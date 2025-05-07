-- name: CreateNewCommunity :one
INSERT INTO community (title, about, thumbnail_url, logo_url, updated_by)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
) RETURNING *;