-- Create a new playlist
-- name: CreatePlaylist :one
INSERT INTO playlist (name, description, thumbnail_url, updated_by)
VALUES (
    $1,  -- Name
    $2,  -- Description
    $3,  -- Thumbnail URL
    $4   -- Updated By
) RETURNING *;

-- Fetch all playlists
-- name: GetAllPlaylists :many
SELECT * FROM playlist;

-- Get a playlist by ID
-- name: GetPlaylistByID :one
SELECT
  id,
  name,
  description,
  thumbnail_url,
  updated_by,
  created_at,
  updated_at
FROM playlist
WHERE id = $1;
