-- name: CreateNewUser :one
INSERT INTO
    "user" (id, user_type, updated_by)
VALUES (
        $1, -- id
        $2, -- UserType
        $3 -- Updated By
    ) RETURNING *;