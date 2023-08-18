-- name: CreateUser :one
INSERT INTO users (
  name,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users 
WHERE name = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
  hashed_password=COALESCE(sqlc.narg(hashed_password), hashed_password),
  full_name=COALESCE(sqlc.narg(full_name), full_name),
  email=COALESCE(sqlc.narg(email), email)
WHERE
  name = sqlc.arg(name)
RETURNING *;

-- name: UpdateUser :one
-- UPDATE users
-- SET
--   hashed_password=CASE
--     WHEN sqlc.arg(set_hashed_password)::boolean = TRUE THEN sqlc.arg(hashed_password)
--     ELSE hashed_password
--   END,
--   full_name=CASE
--     WHEN @set_full_name::boolean = TRUE THEN @full_name
--     ELSE full_name
--   END,
--   email=CASE
--     WHEN @set_email::boolean = TRUE THEN @email
--     ELSE email
--   END
-- WHERE
--   name=@name
-- RETURNING *;