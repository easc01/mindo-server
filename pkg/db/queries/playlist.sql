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

-- unused, left for reference
-- name: GetPlaylistWithTopics :one
SELECT p.id, p.name, p.description, p.code, p.thumbnail_url, p.views, p.created_at, p.updated_at, p.updated_by, COALESCE(
        json_agg (
            t.name
            ORDER BY t.number ASC
        ), '[]'
    ) AS topics
FROM playlist p
    LEFT JOIN topic t ON p.id = t.playlist_id
WHERE
    p.id = $1
GROUP BY
    p.id;

-- name: UpdatePlaylistViewCountById :exec
UPDATE playlist SET views = views + $2 WHERE id = $1;

-- name: GetAllPlaylistsPreviews :many
SELECT
    p.id,
    p.name,
    p.description,
    p.code,
    p.thumbnail_url,
    p.interest_id,
    p.views,
    p.created_at,
    p.updated_at,
    p.updated_by,
    COALESCE(COUNT(t.id), 0) AS topics_count
FROM playlist p
LEFT JOIN topic t ON t.playlist_id = p.id
WHERE $1 = '' OR similarity(p.name, $1) > 0.1
GROUP BY p.id
ORDER BY similarity(p.name, $1) DESC;
