-- name: CreateUserPlaylist :one
INSERT INTO user_playlist (
    user_id,
    playlist_id,
    updated_by
)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, playlist_id) DO UPDATE
SET 
    updated_at = NOW(),
    updated_by = EXCLUDED.updated_by
RETURNING *;
