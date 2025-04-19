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
        $8 -- Updated By
    ) RETURNING *;

-- name: CreateNewUser :one
INSERT INTO
    "user" (id, user_type, updated_by)
VALUES (
        $1, -- id
        $2, -- UserType
        $3 -- Updated By
    ) RETURNING *;

-- name: GetAppUserByUserID :one
SELECT
    u.id AS user_id,
    au.username,
    au.profile_picture_url,
    au.oauth_client_id,
    au.bio,
    au.name,
    au.mobile,
    au.email,
    au.last_login_at,
    au.created_at,
    au.updated_at,
    au.updated_by
FROM app_user au
    JOIN "user" u ON u.id = au.user_id
WHERE
    au.user_id = $1;

-- name: GetAppUserByUsername :one
SELECT
    u.id AS user_id,
    au.username,
    au.profile_picture_url,
    au.oauth_client_id,
    au.bio,
    au.name,
    au.mobile,
    au.email,
    au.last_login_at,
    au.created_at,
    au.updated_at,
    au.updated_by
FROM app_user au
    JOIN "user" u ON u.id = au.user_id
WHERE
    au.username = $1;

-- name: GetAppUserByClientOAuthID :one
SELECT
    u.id AS user_id,
    au.username,
    au.profile_picture_url,
    au.oauth_client_id,
    au.bio,
    au.name,
    au.mobile,
    au.email,
    au.last_login_at,
    au.created_at,
    au.updated_at,
    au.updated_by
FROM app_user au
    JOIN "user" u ON u.id = au.user_id
WHERE
    au.oauth_client_id = $1;

-- name: UpdateUserLastLoginAtByUsername :one
UPDATE app_user
SET
    last_login_at = now()
WHERE
    username = $1 RETURNING user_id,
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

-- name: UpdateUserLastLoginAtByOAuthClientID :one
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