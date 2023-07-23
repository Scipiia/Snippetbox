-- name: CreateSnippet :one
INSERT INTO snippets (
  account_id,
  title,
  content
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetSnippet :one
SELECT * FROM snippets
WHERE id = $1 LIMIT 1;

-- name: ListSnippets :many
SELECT * FROM snippets
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: DeleteSnippet :exec
DELETE FROM snippets
WHERE id = $1;