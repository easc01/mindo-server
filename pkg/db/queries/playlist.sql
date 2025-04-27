-- Create a new playlist
-- name: CreatePlaylist :one
INSERT INTO
    playlist (
        name,
        description,
        thumbnail_url,
        code,
        interest_id,
        updated_by
    )
VALUES (
        $1, -- Name
        $2, -- Description
        $3, -- Thumbnail URL
        $4, -- unique hexcode of playlist
        $5, -- domain/interest id
        $6 -- Updated By
    ) RETURNING *;

-- Get a playlist by ID
-- name: GetPlaylistByID :one
SELECT * FROM playlist WHERE id = $1;