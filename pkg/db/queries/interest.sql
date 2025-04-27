-- name: GetAllInterest :many
SELECT * FROM interest;

-- name: GetInterestByName :one
SELECT * FROM interest WHERE name = $1;