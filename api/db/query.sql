-- name: GetTemple :one
SELECT * FROM temples
WHERE id = $1 LIMIT 1;

-- name: ListTemples :many
SELECT * FROM temples
ORDER BY name_th ASC;

-- name: CreateUser :exec
INSERT INTO users (username, email) VALUES ($1, $2) RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;