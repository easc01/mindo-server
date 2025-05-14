-- name: SaveQuiz :one
INSERT INTO
    "quiz" (
        name,
        thumbnail_url,
        play_count,
        updated_by
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: SaveQuizQuestion :one
INSERT INTO
    "quiz_question" (
        quiz_id,
        question,
        options,
        correct_option,
        updated_by
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetQuestionsByQuizId :many
SELECT * FROM "quiz_question" WHERE quiz_id = $1;