-- name: UpdateUserJoinedCommunityAccess :one
UPDATE user_joined_community
SET
  updated_at = now(),
  updated_by = $3
WHERE community_id = $1 AND user_id = $2
RETURNING *;
