-- name: CreateUser :one
INSERT INTO users (username, avatar, email) VALUES ($1, $2, $3) RETURNING id;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;