-- name: GetFiles :many
SELECT * FROM files

-- name: CreateProfile :one
INSERT INTO profiles (name) VALUES (?) RETURNING *;

-- name: GetConflictById :one
SELECT * FROM conflicts WHERE id = ?
