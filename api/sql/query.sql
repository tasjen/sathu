-- name: GetTemple :one
SELECT * FROM temples
WHERE id = $1 LIMIT 1;

-- name: ListTemples :many
SELECT * FROM temples
ORDER BY thai_name;