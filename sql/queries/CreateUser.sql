-- name: CreateUser :one
INSERT INTO users (
  id,
  created_at,
  updated_at,
  name
) VALUES (
  $1,
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP,
  $2
)
  RETURNING *;
