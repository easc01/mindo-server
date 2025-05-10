-- name: CreateNewAppUser :one
INSERT INTO
    app_user (
        user_id,
        name,
        username,
        email,
        mobile,
        password_hash,
        oauth_client_id,
        color,
        updated_by
    )
VALUES (
        $1, -- id
        $2, -- Name
        $3, -- Username
        $4, -- Email
        $5, -- Mobile
        $6, -- Password Hash
        $7, -- OAuth Client ID
        $8, -- Color
        $9 -- Updated By
    ) RETURNING *;


-- name: UpdateAppUserLastLoginAtByOAuthClientID :one
UPDATE app_user
SET
    last_login_at = now()
WHERE
    oauth_client_id = $1 RETURNING user_id,
    username,
    profile_picture_url,
    bio,
    name,
    mobile,
    email,
    oauth_client_id,
    last_login_at,
    created_at,
    updated_at,
    updated_by;