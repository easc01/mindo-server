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

-- name: UpdateAppUserLastLoginAtByUsername :one
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

-- name: GetAppUserWithInterestsByUserID :one
SELECT 
    au.user_id,
    au.username,
    au.profile_picture_url,
    au.bio,
    au.name,
    au.mobile,
    au.email,
    au.last_login_at,
    au.created_at,
    au.updated_at,
    au.updated_by,
    COALESCE(
        jsonb_agg(
            jsonb_build_object(
                'id', aui.id,
                'name', COALESCE(aui.name, i.name)
            )
        ) FILTER (WHERE aui.id IS NOT NULL),
        '[]'::jsonb
    ) AS interests
FROM 
    app_user au
LEFT JOIN 
    app_user_interest aui ON au.user_id = aui.app_user_id
LEFT JOIN 
    interest i ON aui.interest_id = i.id
WHERE 
    au.user_id = $1
GROUP BY 
    au.user_id;


