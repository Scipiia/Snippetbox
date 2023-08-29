-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
    name,
    email,
    secret_code
)   VALUES (
    $1, $2, $3
)   RETURNING *;