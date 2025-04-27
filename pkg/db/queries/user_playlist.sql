-- name: CreateUserPlaylist :one
INSERT INTO
    user_playlist (
        user_id,
        playlist_id,
        updated_by
    )
VALUES ($1, $2, $3) RETURNING *;