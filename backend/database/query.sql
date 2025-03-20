-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetAllUsers :many
SELECT * FROM users
ORDER BY email;

-- name: CreateUser :one
INSERT INTO users (
  id, name, email
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
set email = ?, name = ?
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;