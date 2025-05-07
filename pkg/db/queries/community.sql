-- name: CreateNewCommunity :one
INSERT INTO community (title, about, thumbnail_url, logo_url, updated_by)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
) RETURNING *;


-- name: CreateNewUserJoinedCommunityById :exec
INSERT INTO user_joined_community (user_id, community_id, updated_by)
VALUES (
  $1,
  $2,
  $3
);


-- name: CreateMessage :one
INSERT INTO "message" (user_id, community_id, content, updated_by)
VALUES (
  $1,
  $2,
  $3,
  $4
) RETURNING *;