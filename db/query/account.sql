-- name: CreateAccount :one
INSERT INTO account (
  login,
  username
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM account
WHERE id = $1 LIMIT 1;

-- name: UpdateAccount :one
UPDATE account
  set username = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM account
WHERE id = $1;