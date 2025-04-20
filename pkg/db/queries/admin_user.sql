-- name: CreateNewAdminUser :one
INSERT INTO
    admin_user (
        user_id,
        name,
        email,
        password_hash,
        updated_by
    )
VALUES (
        $1, -- id
        $2, -- Name
        $3, -- Email
        $4, -- Password Hash
        $5  -- Updated By
    ) RETURNING *;


-- name: GetAdminUserByEmail :one
SELECT
    u.id AS user_id,
    u.user_type,
    au.name,
    au.email,
    au.password_hash,
    au.last_login_at,
    au.created_at,
    au.updated_at,
    au.updated_by
FROM admin_user au
    JOIN "user" u ON u.id = au.user_id
WHERE
    au.email = $1;

-- name: UpdateAdminUserLastLoginByUserId :exec
UPDATE admin_user
SET last_login_at = now()
WHERE user_id = $1;