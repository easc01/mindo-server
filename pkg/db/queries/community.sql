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
WITH inserted_message AS (
  INSERT INTO "message" (user_id, community_id, content, updated_by)
  VALUES ($1, $2, $3, $4)
  RETURNING *
)
SELECT 
  im.*, 
  au.name, 
  au.username, 
  au.color, 
  au.profile_picture_url
FROM inserted_message im
JOIN "app_user" au ON au.user_id = im.user_id;
