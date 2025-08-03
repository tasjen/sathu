-- name: CreateUserWithPassword :one
INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;