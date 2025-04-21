-- name: UpsertUserToken :one
INSERT INTO
    user_token (
        user_id,
        refresh_token,
        role,
        expires_at,
        updated_by
    )
VALUES (
    $1, -- User Id
    $2, -- Refresh Token
    $3, -- Role
    $4, -- Expires At
    $5  -- Updated By
)
ON CONFLICT (user_id)  -- Specify the unique constraint (e.g., user_id)
DO UPDATE SET
    refresh_token = EXCLUDED.refresh_token,
    expires_at = EXCLUDED.expires_at,
    updated_by = EXCLUDED.updated_by
RETURNING *;


-- name: GetUserTokenByRefreshToken :one
SELECT * FROM user_token WHERE refresh_token = $1;