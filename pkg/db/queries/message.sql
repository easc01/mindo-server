-- name: GetMessagePageByCommunityID :many
SELECT m.*, au.username, au.profile_picture_url, au.name
FROM "message" m
JOIN app_user au ON au.user_id = m.user_id
WHERE
    m.community_id = $1
    AND m.created_at < $2
ORDER BY m.created_at DESC
LIMIT 50;